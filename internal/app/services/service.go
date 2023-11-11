package services

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"github.com/DenisKhanov/shorterURL/internal/app/models"
	"github.com/sirupsen/logrus"
	"strings"
)

//go:generate mockgen -source=service.go -destination=mocks/service_mock.go -package=mocks

type Repository interface {
	StoreURLInDB(originalURL, shortURL string) error
	GetShortURLFromDB(originalURL string) (string, error)
	GetOriginalURLFromDB(shortURL string) (string, error)
}
type URLInMemoryRepository interface {
	SaveBatchToFile() error
}
type URLInDBRepository interface {
	GetShortBatchURLFromDB(batchURLRequests []models.URLRequest) (map[string]string, error)
	StoreBatchURLInDB(batchURLtoStores map[string]string) error
}
type Encoder interface {
	CryptoBase62Encode() string
}

type Services struct {
	repository            Repository
	URLInDBRepository     URLInDBRepository
	URLInMemoryRepository URLInMemoryRepository
	encoder               Encoder
	baseURL               string
}

func NewServices(repository Repository, URLInDBRepository URLInDBRepository, URLInMemoryRepository URLInMemoryRepository, encoder Encoder, baseURL string) *Services {
	return &Services{
		repository:            repository,
		URLInDBRepository:     URLInDBRepository,
		URLInMemoryRepository: URLInMemoryRepository,
		encoder:               encoder,
		baseURL:               baseURL,
	}
}

func (s Services) GetBatchJSONShortURL(batchURLRequests []models.URLRequest) ([]models.URLResponse, error) {
	fmt.Println("Service GetBatchJSONShortURL run")
	var batchURLtoStores = make(map[string]string, len(batchURLRequests))
	var batchURLResponses []models.URLResponse
	shortsURL, err := s.URLInDBRepository.GetShortBatchURLFromDB(batchURLRequests)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	for _, value := range batchURLRequests {
		if shortURL, ok := shortsURL[value.OriginalURL]; ok {
			batchURLResponses = append(batchURLResponses, models.URLResponse{CorrelationID: value.CorrelationID, ShortURL: s.baseURL + "/" + shortURL})
		} else {
			shortURL = s.encoder.CryptoBase62Encode()
			batchURLtoStores[shortURL] = value.OriginalURL
			batchURLResponses = append(batchURLResponses, models.URLResponse{CorrelationID: value.CorrelationID, ShortURL: s.baseURL + "/" + shortURL})
		}
	}
	err = s.URLInDBRepository.StoreBatchURLInDB(batchURLtoStores)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return batchURLResponses, nil
}

// GetShortURL returns the short URL
func (s Services) GetShortURL(url string) (string, error) {
	shortURL, err := s.repository.GetShortURLFromDB(url)
	if err != nil {
		shortURL = s.encoder.CryptoBase62Encode()
		err = s.repository.StoreURLInDB(url, shortURL)
		if err != nil {
			return "", err
		}
	} else {
		return s.baseURL + "/" + shortURL, nil
	}
	return s.baseURL + "/" + shortURL, nil
}

// GetOriginalURL returns the origin URL for the given short URL
func (s Services) GetOriginalURL(shortURL string) (string, error) {
	originURL, err := s.repository.GetOriginalURLFromDB(shortURL)
	if err != nil {
		return "", err
	}
	return originURL, nil
}

// CryptoBase62Encode generates a unique string that is a
// Base62-encoded representation of a 42-bit random number.
// The random number is generated using a cryptographically
// secure random number generator.
// The returned string has a length of up to 7 characters
func (s Services) CryptoBase62Encode() string {
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
