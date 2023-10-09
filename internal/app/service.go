package app

import (
	"fmt"
	"strings"
)

//go:generate mockgen -source=service.go -destination=mocks/service_mock.go -package=mocks
type Service interface {
	GetShortURL(url string) (string, error)
	GetOriginURL(shortURL string) (string, error)
}

type Services struct {
	storage Repository
	service Service
}

func NewServices(storage Repository, service Service) *Services {
	return &Services{
		storage: storage,
		service: service,
	}
}

// GetShortURL returns the short URL ("http://localhost:8080/"+shortURL)
func (s Services) GetShortURL(url string) (string, error) {
	value, exists := s.storage.GetShortURL(url)
	if exists {
		return "http://localhost:8080/" + value, nil
	} else {
		shortURL := base62Encode(s.storage.GetID())
		err := s.storage.StoreURL(url, shortURL)
		if err != nil {
			return "", err
		}
		return "http://localhost:8080/" + shortURL, nil
	}
}

// GetOriginURL returns the origin URL for the given short URL
func (s Services) GetOriginURL(shortURL string) (string, error) {
	originURL, exists := s.storage.GetOriginalURL(shortURL)
	if !exists {
		return "", fmt.Errorf("http://localhost:8080/%s not found", shortURL)
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
