// Package repositories provides implementations for interacting with the data storage.
// It includes functionality to store, retrieve, and manage shortened URLs in a PostgreSQL database.
package repositories

import (
	"context"
	"errors"
	"fmt"
	"github.com/DenisKhanov/shorterURL/internal/app/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

// URLInDBRepo represents the repository for storing and retrieving URLs in the PostgreSQL database.
type URLInDBRepo struct {
	UserID uuid.UUID     `json:"id"`
	DB     *pgxpool.Pool //opened in main func DB pool connections
}

// NewURLInDBRepo creates a new instance of URLInDBRepo with the provided database connection pool.
// It also initializes the database table for storing shortened URLs.
func NewURLInDBRepo(DB *pgxpool.Pool) *URLInDBRepo {
	storage := &URLInDBRepo{
		UserID: uuid.Nil,
		DB:     DB,
	}
	storage.CreateBDTable()
	return storage
}

// CreateBDTable creates the "shorted_URL" table in the database if it doesn't already exist.
func (d *URLInDBRepo) CreateBDTable() error {
	ctx := context.Background()
	sqlQuery := `
		CREATE TABLE IF NOT EXISTS shorted_URL (
		user_id UUID NOT NULL,
		short_url VARCHAR(250) NOT NULL,
		original_url VARCHAR(4096) NOT NULL UNIQUE,
		deleted_flag bool NOT NULL DEFAULT false
	)`
	_, err := d.DB.Exec(ctx, sqlQuery)
	if err != nil {
		logrus.Errorf("don't create table shorted_URL: %v", err)
		return err
	}
	logrus.Info("Successfully created table shorted_URL")
	return nil
}

// StoreURLInDB saves a mapping between an original URL and its shortened version in the database.
// It returns an error if the saving process fails.
func (d *URLInDBRepo) StoreURLInDB(ctx context.Context, originalURL, shortURL string) error {
	userID, ok := ctx.Value(models.UserIDKey).(uuid.UUID)
	if !ok {
		logrus.Errorf("context value is not userID: %v", userID)
	}
	const sqlQuery = `INSERT INTO shorted_URL (user_id, original_url, short_url) VALUES ($1, $2, $3) ON CONFLICT (original_url) DO NOTHING`
	_, err := d.DB.Exec(ctx, sqlQuery, userID, originalURL, shortURL)
	if err != nil {
		logrus.Error("url don't save in database ", err)
		return err
	}
	return nil
}

// StoreBatchURLInDB saves multiple URL mappings in the database in a batch operation.
// The input is a map where keys are shortened URLs and values are the corresponding original URLs.
// It returns an error if the batch saving process fails.
func (d *URLInDBRepo) StoreBatchURLInDB(ctx context.Context, batchURLtoStores map[string]string) error {
	userID, ok := ctx.Value(models.UserIDKey).(uuid.UUID)
	if !ok {
		logrus.Errorf("context value is not userID: %v", userID)
	}
	tx, err := d.DB.Begin(ctx)
	if err != nil {
		return err
	}
	const sqlQuery = `INSERT INTO shorted_URL (user_id, original_url, short_url) VALUES ($1, $2, $3) ON CONFLICT (original_url) DO NOTHING`
	_, err = tx.Prepare(ctx, "store_batch_url", sqlQuery)
	if err != nil {
		return err
	}
	for shortURL, originalURL := range batchURLtoStores {
		_, err = tx.Exec(ctx, "store_batch_url", userID, originalURL, shortURL)
		if err != nil {
			logrus.Error("url don't save in database ", err)
			tx.Rollback(ctx)
			return err
		}
	}
	return tx.Commit(ctx)
}

// MarkURLsAsDeleted marks user URLs as deleted in DB
func (d *URLInDBRepo) MarkURLsAsDeleted(ctx context.Context, URLSToDel []string) error {
	if len(URLSToDel) == 0 {
		return nil
	}
	userID, ok := ctx.Value(models.UserIDKey).(uuid.UUID)
	if !ok {
		logrus.Errorf("context value is not userID: %v", userID)
		return fmt.Errorf("invalid user context")
	}
	tx, err := d.DB.Begin(ctx)
	if err != nil {
		logrus.Error("Failed to begin transaction: ", err)
		return err
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				logrus.Errorf("Failed to rollback transaction: %v", rollbackErr)
			}
		}
	}()

	const sqlQuery = `UPDATE shorted_URL SET deleted_flag = true WHERE short_url = ANY($1) AND user_id = $2`
	_, err = tx.Exec(ctx, sqlQuery, URLSToDel, userID)
	if err != nil {
		logrus.Error("Failed to mark URLs as deleted: ", err)
		return err
	}
	logrus.Infof("Complete mark URLs as deleted: %s, %d", URLSToDel, userID)
	return tx.Commit(ctx)
}

// GetOriginalURLFromDB retrieves the original URL corresponding to a given shortened URL from the database.
// It returns the original URL and any error encountered during the retrieval.
func (d *URLInDBRepo) GetOriginalURLFromDB(ctx context.Context, shortURL string) (string, error) {
	const selectQuery = `SELECT original_url, deleted_flag FROM shorted_URL WHERE short_url = $1`
	var originalURL string
	var deletedFlag bool
	err := d.DB.QueryRow(ctx, selectQuery, shortURL).Scan(&originalURL, &deletedFlag)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("original URL not found")
		}
		logrus.Error("error querying for short URL: ", err)

		return "", fmt.Errorf("error querying for short URL: %w", err)
	}
	if deletedFlag {
		return "", models.ErrURLDeleted
	}
	return originalURL, nil
}

// GetShortURLFromDB retrieves the shortened version of a given original URL from the database.
// It returns the shortened URL and any error encountered during the retrieval.
func (d *URLInDBRepo) GetShortURLFromDB(ctx context.Context, originalURL string) (string, error) {
	const selectQuery = `SELECT short_url FROM shorted_URL WHERE original_url = $1`
	var shortURL string
	err := d.DB.QueryRow(ctx, selectQuery, originalURL).Scan(&shortURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("short URL not found: %w", err)
		}
		logrus.Error("error querying for original URL: ", err)

		return "", fmt.Errorf("error querying for original URL: %w", err)
	}
	return shortURL, nil
}

// GetUserURLSFromDB takes a slice of models.URL objects for a specific user from DB
func (d *URLInDBRepo) GetUserURLSFromDB(ctx context.Context) ([]models.URL, error) {
	const selectQuery = `SELECT short_url,original_url FROM shorted_URL WHERE user_id = $1`
	userID, ok := ctx.Value(models.UserIDKey).(uuid.UUID)
	if !ok {
		logrus.Errorf("context value is not userID: %v", userID)
	}
	rows, err := d.DB.Query(ctx, selectQuery, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("userID not found: %w", err)
		}
		logrus.Error("error querying for user usersURLS: ", err)
		return nil, fmt.Errorf("error querying for user usersURLS: %w", err)
	}
	defer rows.Close()

	var userURLS []models.URL
	for rows.Next() {
		rowResult := models.URL{}
		if err = rows.Scan(&rowResult.ShortURL, &rowResult.OriginalURL); err != nil {
			logrus.Error(err)
		}
		userURLS = append(userURLS, rowResult)
	}
	if err = rows.Err(); err != nil {
		logrus.Error(err)
		return nil, err
	}
	return userURLS, nil
}

// GetShortBatchURLFromDB retrieves multiple shortened URLs corresponding to a batch of original URLs from the database.
// The input is a slice of URLRequest objects containing original URLs.
//
//	It returns found in database a map of original URLs to their shortened counterparts and any error encountered during the retrieval.
func (d *URLInDBRepo) GetShortBatchURLFromDB(ctx context.Context, batchURLRequests []models.URLRequest) (map[string]string, error) {
	var shortsURL = make(map[string]string, len(batchURLRequests))
	var shortURL string
	if d.DB == nil {
		fmt.Println("Repository not pool")
	}
	tx, err := d.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	const selectQuery = `SELECT short_url FROM shorted_URL WHERE original_url = $1`
	for _, request := range batchURLRequests {
		err = tx.QueryRow(ctx, selectQuery, request.OriginalURL).Scan(&shortURL)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				continue
			}
			tx.Rollback(ctx)
			logrus.Error("error querying for original URL: ", err)
			return nil, fmt.Errorf("error querying for original URL: %w", err)
		}
		shortsURL[request.OriginalURL] = shortURL
	}

	return shortsURL, tx.Commit(ctx)
}

// Stats retrieves the statistics of URLs and users from the database.
//
// This method executes a SQL query to retrieve the count of URLs and users from the 'shorted_url' table.
// It then constructs a Stats struct containing the counts and returns it along with any error encountered.
func (d *URLInDBRepo) Stats(ctx context.Context) (models.Stats, error) {
	const selectQuery = `SELECT 
    					 COUNT(short_url),
    					 COUNT(DISTINCT user_id)
						 FROM shorted_url`
	var urls, users int
	err := d.DB.QueryRow(ctx, selectQuery).Scan(&urls, &users)
	if err != nil {
		logrus.Error("error querying for count urls or users: ", err)
		return models.Stats{}, fmt.Errorf("error querying for count urls or users: %w", err)
	}
	stats := models.Stats{Urls: urls, Users: users}
	return stats, nil
}
