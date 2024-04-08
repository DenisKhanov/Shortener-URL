// Package services provides the business logic for managing shortened URLs.
// It includes functionality to generate, store, retrieve, and delete URLs.
package services

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"github.com/DenisKhanov/shorterURL/internal/models"
	"github.com/DenisKhanov/shorterURL/internal/repositories"
	"github.com/sirupsen/logrus"
	"net/url"
	"time"
)

// Repository defines the interface for interacting with the storage backend.
//
//go:generate mockgen -source=service.go -destination=mocks/service_mock.go -package=mocks
type Repository interface {
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
	// GetUserURLS takes a slice of models.URL objects for a specific user from DB
	GetUserURLS(ctx context.Context) ([]models.URL, error)
	// MarkURLsAsDeleted marks user URLs as deleted in DB
	MarkURLsAsDeleted(ctx context.Context, URLSToDel []string) error
	// GetStats retrieves the statistics of URLs and users from the database.
	GetStats(ctx context.Context) (models.Stats, error)
}

var _ Repository = (*repositories.URLInMemoryRepo)(nil)
var _ Repository = (*repositories.URLInDBRepo)(nil)

// Encoder defines the interface for encoding unique short URLs.
type Encoder interface {
	CryptoBase62Encode() string
}

// ShortURLServices represents the service for managing shortened URLs.
type ShortURLServices struct {
	repository Repository
	encoder    Encoder
	baseURL    string
}

// URLInMemoryRepository defines the interface for an in-memory repository to save batch data to a file.
type URLInMemoryRepository interface {
	SaveBatchToFile() error
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

// finalURLBuilder the function combines the base url and the shortened url into a single link
func (s ShortURLServices) finalURLBuilder(shortURL string) string {
	resultURL, err := url.JoinPath(s.baseURL, shortURL)
	if err != nil {
		logrus.Error(err)
	}
	return resultURL
}

// GetBatchShortURL takes a slice of models.URLRequest objects, each containing a URL to be shortened,
// and returns a slice of models.URLResponse objects, each containing the original and shortened URL.
// This method is intended for processing multiple URLs at once, improving efficiency for bulk operations.
// Returns an error if any of the URLs cannot be processed or if an internal error occurs.
func (s ShortURLServices) GetBatchShortURL(ctx context.Context, batchURLRequests []models.URLRequest) ([]models.URLResponse, error) {
	shortsURL, err := s.repository.GetShortBatchURL(ctx, batchURLRequests)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	var batchURLtoStores = make(map[string]string, len(batchURLRequests))
	var batchURLResponses []models.URLResponse
	for _, value := range batchURLRequests {
		if shortURL, ok := shortsURL[value.OriginalURL]; ok {
			batchURLResponses = append(batchURLResponses, models.URLResponse{CorrelationID: value.CorrelationID, ShortURL: s.finalURLBuilder(shortURL)})
		} else {
			shortURL = s.encoder.CryptoBase62Encode()
			batchURLtoStores[shortURL] = value.OriginalURL
			batchURLResponses = append(batchURLResponses, models.URLResponse{CorrelationID: value.CorrelationID, ShortURL: s.finalURLBuilder(shortURL)})
		}
	}

	err = s.repository.StoreBatchURL(ctx, batchURLtoStores)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return batchURLResponses, nil
}

// GetShortURL takes original URL and returns its shortened version.
// If the URL has already been shortened, it returns the existing shortened URL.
// If the URL is new, it generates a new shortened URL.
// Returns an error if the URL cannot be shortened or if any internal error occurs.
func (s ShortURLServices) GetShortURL(ctx context.Context, URL string) (string, error) {
	shortURL, err := s.repository.GetShortURL(ctx, URL)
	if err != nil {
		shortURL = s.encoder.CryptoBase62Encode()
		err = s.repository.StoreURL(ctx, URL, shortURL)
		if err != nil {
			return "", err
		}
		return s.finalURLBuilder(shortURL), nil
	}
	return s.finalURLBuilder(shortURL), models.ErrURLFound
}

// GetOriginalURL takes a shortened URL and returns the original URL it points to.
// If the shortened URL does not exist or is invalid, an error is returned.
// Useful for redirecting shortened URLs to their original destinations.
func (s ShortURLServices) GetOriginalURL(ctx context.Context, shortURL string) (string, error) {
	originURL, err := s.repository.GetOriginalURL(ctx, shortURL)
	if err != nil {
		return "", err
	}
	return originURL, nil
}

// GetUserURLS takes a slice of models.URL objects for a specific user
func (s ShortURLServices) GetUserURLS(ctx context.Context) ([]models.URL, error) {
	userURLS, err := s.repository.GetUserURLS(ctx)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	fullShortUserURLS := make([]models.URL, len(userURLS))
	for i, v := range userURLS {
		fullShortUserURLS[i].ShortURL = s.finalURLBuilder(v.ShortURL)
		fullShortUserURLS[i].OriginalURL = v.OriginalURL
	}
	return fullShortUserURLS, nil
}

// AsyncDeleteUserURLs async runs requests to DB for mark user URLs as deleted
func (s ShortURLServices) AsyncDeleteUserURLs(ctx context.Context, URLSToDel []string) {
	go func() {
		asyncCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), time.Minute)
		defer cancel()
		if err := s.repository.MarkURLsAsDeleted(asyncCtx, URLSToDel); err != nil {
			logrus.Error(err)
		}
	}()
}

// ServiceStats retrieves the statistics of URLs and users from the service's repository.
//
// This method delegates the retrieval of statistics to the repository's GetStats method.
// It then returns the obtained statistics and any error encountered during the retrieval process.
func (s ShortURLServices) ServiceStats(ctx context.Context) (models.Stats, error) {
	stats, err := s.repository.GetStats(ctx)
	if err != nil {
		return models.Stats{}, err
	}
	return stats, err
}

// CryptoBase62Encode generates a unique string that is a
// Base62-encoded representation of a 42-bit random number.
// The random number is generated using a cryptographically
// secure random number generator.
// The returned string has a length of up to 7 characters
func (s ShortURLServices) CryptoBase62Encode() string {
	b := make([]byte, 8) // uint64 состоит из 8 байт, но мы будем использовать только 42 бита
	_, _ = rand.Read(b)
	num := binary.BigEndian.Uint64(b) & ((1 << 42) - 1) // Обнуление всех бит, кроме младших 42 бит
	chars := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	var shortURL = make([]byte, 0, 8)
	for num > 0 {
		remainder := num % 62
		shortURL = append(shortURL, chars[remainder])
		num = num / 62
	}
	return string(shortURL)
}
