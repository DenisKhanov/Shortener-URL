// Package url provides the business logic for managing shortened URLs.
// It includes functionality to generate, store, retrieve, and delete URLs.
package url

import (
	"context"
	"github.com/DenisKhanov/shorterURL/internal/models"
	url2 "github.com/DenisKhanov/shorterURL/internal/repositories/url"
	"github.com/sirupsen/logrus"
	"net/url"
)

// Repository defines the interface for interacting with the storage backend.
//
//go:generate mockgen -source=service.go -destination=mocks/service_mock.go -package=mocks
type Repository interface {
	// Ping checks the database connection or repository created.
	Ping(ctx context.Context) error
	// StoreURL saves a mapping between an original URL and its shortened version in the database.
	// It returns an error if the saving process fails.
	StoreURL(ctx context.Context, originalURL, shortURL string) error
	// GetShortURL retrieves the shortened version of a given original URL from the database.
	// It returns the shortened URL and any error encountered during the retrieval.
	GetShortURL(ctx context.Context, originalURL string) (string, error)
	// GetOriginalURL retrieves the original URL corresponding to a given shortened URL from the database.
	// It returns the original URL and any error encountered during the retrieval.
	GetOriginalURL(ctx context.Context, shortURL string) (string, error)
	// StoreBatchURL saves multiple URL mappings in the database in a batch operation.
	// The input is a map where keys are shortened URLs and values are the corresponding original URLs.
	// It returns an error if the batch saving process fails.
	StoreBatchURL(ctx context.Context, batchURLtoStores map[string]string) error
	// GetShortBatchURL retrieves multiple shortened URLs corresponding to a batch of original URLs from the database.
	// The input is a slice of URLRequest objects containing original URLs.
	//  It returns found in database a map of original URLs to their shortened counterparts and any error encountered during the retrieval.
	GetShortBatchURL(ctx context.Context, batchURLRequests []models.URLRequest) (map[string]string, error)
	// GetUserURLs takes a slice of models.URL objects for a specific user from DB
	GetUserURLs(ctx context.Context) ([]models.URL, error)
	// MarkURLsAsDeleted marks user URLs as deleted in DB
	MarkURLsAsDeleted(ctx context.Context, URLSToDel []string) error
	// GetStats retrieves the statistics of URLs and users from the database.
	GetStats(ctx context.Context) (models.Stats, error)
}

// InMemoryRepository defines the interface for an in-memory repository to save batch data to a file.
type InMemoryRepository interface {
	SaveBatchToFile() error
}

// Encoder defines the interface for encoding unique short URLs.
type Encoder interface {
	CryptoBase62Encode() string
}

// checking interface compliance at the compiler level
var _ Repository = (*url2.URLInMemoryRepo)(nil)
var _ Repository = (*url2.URLInDBRepo)(nil)

// ShortURLServices represents the service for managing shortened URLs.
type ShortURLServices struct {
	repository Repository
	save       InMemoryRepository
	encoder    Encoder
	baseURL    string
}

// NewShortURLServices creates a new instance of ShortURLServices.
// It takes a repository for data storage, an encoder for generating short URLs, and a base URL.
func NewShortURLServices(repository Repository, encoder Encoder, baseURL string) *ShortURLServices {
	parsedBaseURL, err := url.Parse(baseURL)
	if err != nil {
		logrus.Error(err)
	}
	return &ShortURLServices{
		repository: repository,
		encoder:    encoder,
		baseURL:    parsedBaseURL.String(),
	}
}
