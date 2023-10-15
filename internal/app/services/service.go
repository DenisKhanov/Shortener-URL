package services

import (
	"fmt"
	"strings"
)

//go:generate mockgen -source=service.go -destination=mocks/service_mock.go -package=mocks
type Repository interface {
	StoreURLSInDB(originalURL, shortURL string) error
	GetShortURLFromDB(originalURL string) (string, bool)
	GetOriginalURLFromDB(shortURL string) (string, bool)
	GetIDFromDB() int
}
type Services struct {
	repository Repository
	baseURL    string
}

func NewServices(storage Repository, baseURL string) *Services {
	return &Services{
		repository: storage,
		baseURL:    baseURL,
	}
}

// GetShortURL returns the short URL
func (s Services) GetShortURL(url string) (string, error) {
	value, exists := s.repository.GetShortURLFromDB(url)
	if exists {
		return s.baseURL + "/" + value, nil
	} else {
		shortURL := base62Encode(s.repository.GetIDFromDB())
		err := s.repository.StoreURLSInDB(url, shortURL)
		if err != nil {
			return "", err
		}
		return s.baseURL + "/" + shortURL, nil
	}
}

// GetOriginalURL returns the origin URL for the given short URL
func (s Services) GetOriginalURL(shortURL string) (string, error) {
	originURL, exists := s.repository.GetOriginalURLFromDB(shortURL)
	if !exists {
		return "", fmt.Errorf("%s/%s not found", s.baseURL, shortURL)
	}
	return originURL, nil

}

// base62Encode returns the base 64 encoded string
func base62Encode(n int) string {
	chars := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	if n == 0 {
		return string(chars[0])
	}
	var shortURL strings.Builder
	for n > 0 {
		shortURL.WriteString(string(chars[n%62]))
		n = n / 62
	}
	return shortURL.String()
}
