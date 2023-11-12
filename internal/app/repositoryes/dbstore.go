package repositoryes

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type URLInDBRepo struct {
	ID          uint8  `json:"id"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	DB          *pgxpool.Pool
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
func (d *URLInDBRepo) StoreURLSInDB(originalURL, shortURL string) error {
	ctx := context.Background()
	const sqlQuery = `INSERT INTO shortedurl (originalurl, shorturl) VALUES ($1, $2) ON CONFLICT (shorturl) DO NOTHING`
	_, err := d.DB.Exec(ctx, sqlQuery, originalURL, shortURL)
	if err != nil {
		logrus.Error("url don't save in database ", err)
		return err
	}
	return nil
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
