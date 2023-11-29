package services

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"github.com/DenisKhanov/shorterURL/internal/app/models"
	"github.com/sirupsen/logrus"
	"net/url"
	"strings"
	"time"
)

// Repository defines the interface for interacting with the storage backend.
//
//go:generate mockgen -source=service.go -destination=mocks/service_mock.go -package=mocks
type Repository interface {
	// StoreURLInDB saves a mapping between an original URL and its shortened version in the database.
	// It returns an error if the saving process fails.
	StoreURLInDB(ctx context.Context, originalURL, shortURL string) error
	// GetShortURLFromDB retrieves the shortened version of a given original URL from the database.
	// It returns the shortened URL and any error encountered during the retrieval.
	GetShortURLFromDB(ctx context.Context, originalURL string) (string, error)
	// GetOriginalURLFromDB retrieves the original URL corresponding to a given shortened URL from the database.
	// It returns the original URL and any error encountered during the retrieval.
	GetOriginalURLFromDB(ctx context.Context, shortURL string) (string, error)
	// StoreBatchURLInDB saves multiple URL mappings in the database in a batch operation.
	// The input is a map where keys are shortened URLs and values are the corresponding original URLs.
	// It returns an error if the batch saving process fails.
	StoreBatchURLInDB(ctx context.Context, batchURLtoStores map[string]string) error
	// GetShortBatchURLFromDB retrieves multiple shortened URLs corresponding to a batch of original URLs from the database.
	// The input is a slice of URLRequest objects containing original URLs.
	//  It returns found in database a map of original URLs to their shortened counterparts and any error encountered during the retrieval.
	GetShortBatchURLFromDB(ctx context.Context, batchURLRequests []models.URLRequest) (map[string]string, error)
	// GetUserURLSFromDB takes a slice of models.URL objects for a specific user from DB
	GetUserURLSFromDB(ctx context.Context) ([]models.URL, error)
	MarkURLsAsDeleted(ctx context.Context, URLSToDel []string) error
}
type Encoder interface {
	CryptoBase62Encode() string
}

type ShortURLServices struct {
	repository Repository
	encoder    Encoder
	baseURL    string
}
type URLInMemoryRepository interface {
	SaveBatchToFile() error
}

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
func (s ShortURLServices) GetBatchShortURL(ctx context.Context, batchURLRequests []models.URLRequest) ([]models.URLResponse, error) {
	shortsURL, err := s.repository.GetShortBatchURLFromDB(ctx, batchURLRequests)
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

	err = s.repository.StoreBatchURLInDB(ctx, batchURLtoStores)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return batchURLResponses, nil
}
func (s ShortURLServices) GetShortURL(ctx context.Context, URL string) (string, error) {
	shortURL, err := s.repository.GetShortURLFromDB(ctx, URL)
	if err != nil {
		shortURL = s.encoder.CryptoBase62Encode()
		err = s.repository.StoreURLInDB(ctx, URL, shortURL)
		if err != nil {
			return "", err
		}
		return s.finalURLBuilder(shortURL), nil
	}
	return s.finalURLBuilder(shortURL), models.ErrURLFound
}
func (s ShortURLServices) GetOriginalURL(ctx context.Context, shortURL string) (string, error) {
	originURL, err := s.repository.GetOriginalURLFromDB(ctx, shortURL)
	if err != nil {
		return "", err
	}
	return originURL, nil
}
func (s ShortURLServices) GetUserURLS(ctx context.Context) ([]models.URL, error) {
	userURLS, err := s.repository.GetUserURLSFromDB(ctx)
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
func (s ShortURLServices) AsyncDeleteUserURLs(ctx context.Context, URLSToDel []string) {
	go func() {
		asyncCtx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		userID, ok := ctx.Value(models.UserIDKey).(uint32)
		if !ok {
			logrus.Errorf("context value is not userID: %v", userID)
		}
		asyncCtx = context.WithValue(asyncCtx, models.UserIDKey, userID)
		defer func() {
			if r := recover(); r != nil {
				logrus.Errorf("Recovered in AsyncDeleteUserURLs: %v", r)
			}
		}()
		if err := s.repository.MarkURLsAsDeleted(asyncCtx, URLSToDel); err != nil {
			logrus.Error(err)
		}
	}()
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
	var shortURL strings.Builder
	for num > 0 {
		remainder := num % 62
		shortURL.WriteString(string(chars[remainder]))
		num = num / 62
	}
	return shortURL.String()
}
