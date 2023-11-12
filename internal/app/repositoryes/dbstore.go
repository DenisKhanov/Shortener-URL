package repositoryes

import (
	"context"
	"errors"
	"fmt"
	"github.com/DenisKhanov/shorterURL/internal/app/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

// URLInDBRepo auxiliary structure for deserialization from jSON for save to database
type URLInDBRepo struct {
	ID          uint8         `json:"id"`
	ShortURL    string        `json:"short_url"`
	OriginalURL string        `json:"original_url"`
	DB          *pgxpool.Pool //opened in main func DB pool connections
}

func NewURLInDBRepo(DB *pgxpool.Pool) *URLInDBRepo {
	storage := &URLInDBRepo{
		ID:          0,
		ShortURL:    "",
		OriginalURL: "",
		DB:          DB,
	}
	storage.CreateBDTable()
	return storage
}

func (d *URLInDBRepo) CreateBDTable() {
	ctx := context.Background()
	sqlQuery := `
		CREATE TABLE IF NOT EXISTS shortedurl (
		"id" SERIAL PRIMARY KEY,
		"shorturl" VARCHAR(250) NOT NULL UNIQUE,
		"originalurl" VARCHAR(4096) NOT NULL
	)`
	_, err := d.DB.Exec(ctx, sqlQuery)
	if err != nil {
		logrus.Panicf("don't create table shortedurl: %v", err)
	}
	logrus.Info("Successfully created table shortedurl")
}
func (d *URLInDBRepo) StoreURLInDB(originalURL, shortURL string) error {
	ctx := context.Background()
	const sqlQuery = `INSERT INTO shortedurl (originalurl, shorturl) VALUES ($1, $2) ON CONFLICT (shorturl) DO NOTHING`
	_, err := d.DB.Exec(ctx, sqlQuery, originalURL, shortURL)
	if err != nil {
		logrus.Error("url don't save in database ", err)
		return err
	}
	return nil
}
func (d *URLInDBRepo) StoreBatchURLInDB(batchURLtoStores map[string]string) error {
	ctx := context.Background()
	tx, err := d.DB.Begin(ctx)
	if err != nil {
		return err
	}
	fmt.Println(2)
	const sqlQuery = `INSERT INTO shortedurl (originalurl, shorturl) VALUES ($1, $2) ON CONFLICT (shorturl) DO NOTHING`
	_, err = tx.Prepare(ctx, "store_batch_url", sqlQuery)
	if err != nil {
		return err
	}
	for shortURL, originalURL := range batchURLtoStores {
		_, err = tx.Exec(ctx, "store_batch_url", originalURL, shortURL)
		if err != nil {
			logrus.Error("url don't save in database ", err)
			tx.Rollback(ctx)
			return err
		}
	}
	return tx.Commit(ctx)
}
func (d *URLInDBRepo) GetOriginalURLFromDB(shortURL string) (string, error) {
	ctx := context.Background()
	const selectQuery = `SELECT originalurl FROM shortedurl WHERE shorturl = $1`
	var originalURL string
	err := d.DB.QueryRow(ctx, selectQuery, shortURL).Scan(&originalURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("original URL not found")
		}
		logrus.Error("error querying for short URL: ", err)

		return "", fmt.Errorf("error querying for short URL: %w", err)
	}
	return originalURL, nil
}
func (d *URLInDBRepo) GetShortURLFromDB(originalURL string) (string, error) {
	ctx := context.Background()

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
func (d *URLInDBRepo) GetShortBatchURLFromDB(batchURLRequests []models.URLRequest) (map[string]string, error) {
	var shortsURL = make(map[string]string, len(batchURLRequests))
	var shortURL string
	fmt.Println("Repository GetShortBatchURLFromDB run")
	ctx := context.Background()
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
