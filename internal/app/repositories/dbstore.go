package repositories

import (
	"context"
	"errors"
	"fmt"
	"github.com/DenisKhanov/shorterURL/internal/app/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type URLInDBRepo struct {
	UserID uint8         `json:"id"`
	DB     *pgxpool.Pool //opened in main func DB pool connections
}

func NewURLInDBRepo(DB *pgxpool.Pool) *URLInDBRepo {
	storage := &URLInDBRepo{
		UserID: 0,
		DB:     DB,
	}
	storage.CreateBDTable()
	return storage
}

func (d *URLInDBRepo) CreateBDTable() error {
	ctx := context.Background()
	sqlQuery := `
		CREATE TABLE IF NOT EXISTS shortedurl (
		"userid" integer NOT NULL,
		"shorturl" VARCHAR(250) NOT NULL,
		"originalurl" VARCHAR(4096) NOT NULL UNIQUE,
		"deletedflag" bool NOT NULL DEFAULT false
	)`
	_, err := d.DB.Exec(ctx, sqlQuery)
	if err != nil {
		logrus.Errorf("don't create table shortedurl: %v", err)
		return err
	}
	logrus.Info("Successfully created table shortedurl")
	return nil
}
func (d *URLInDBRepo) StoreURLInDB(ctx context.Context, originalURL, shortURL string) error {
	userID, ok := ctx.Value(models.UserIDKey).(uint32)
	if !ok {
		logrus.Errorf("context value is not userID: %v", userID)
	}
	const sqlQuery = `INSERT INTO shortedurl (userid, originalurl, shorturl) VALUES ($1, $2, $3) ON CONFLICT (originalurl) DO NOTHING`
	_, err := d.DB.Exec(ctx, sqlQuery, userID, originalURL, shortURL)
	if err != nil {
		logrus.Error("url don't save in database ", err)
		return err
	}
	return nil
}
func (d *URLInDBRepo) StoreBatchURLInDB(ctx context.Context, batchURLtoStores map[string]string) error {
	userID, ok := ctx.Value(models.UserIDKey).(uint32)
	if !ok {
		logrus.Errorf("context value is not userID: %v", userID)
	}
	tx, err := d.DB.Begin(ctx)
	if err != nil {
		return err
	}
	const sqlQuery = `INSERT INTO shortedurl (userid, originalurl, shorturl) VALUES ($1, $2, $3) ON CONFLICT (originalurl) DO NOTHING`
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
func (d *URLInDBRepo) MarkURLsAsDeleted(ctx context.Context, URLSToDel []string) error {
	if len(URLSToDel) == 0 {
		return nil
	}
	userID, ok := ctx.Value(models.UserIDKey).(uint32)
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

	const sqlQuery = `UPDATE shortedurl SET deletedflag = true WHERE shorturl = ANY($1) AND userid = $2`
	_, err = tx.Exec(ctx, sqlQuery, URLSToDel, userID)
	if err != nil {
		logrus.Error("Failed to mark URLs as deleted: ", err)
		return err
	}
	logrus.Infof("Complete mark URLs as deleted: %s, %d", URLSToDel, userID)
	return tx.Commit(ctx)
}
func (d *URLInDBRepo) GetOriginalURLFromDB(ctx context.Context, shortURL string) (string, error) {
	const selectQuery = `SELECT originalurl, deletedflag FROM shortedurl WHERE shorturl = $1`
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
func (d *URLInDBRepo) GetShortURLFromDB(ctx context.Context, originalURL string) (string, error) {
	const selectQuery = `SELECT shorturl FROM shortedurl WHERE originalurl = $1`
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
func (d *URLInDBRepo) GetUserURLSFromDB(ctx context.Context) ([]models.URL, error) {
	const selectQuery = `SELECT shorturl,originalurl FROM shortedurl WHERE userid = $1`
	userID, ok := ctx.Value(models.UserIDKey).(uint32)
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
	const selectQuery = `SELECT shorturl FROM shortedurl WHERE originalurl = $1`
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
