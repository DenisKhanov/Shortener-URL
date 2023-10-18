package services

import (
	"crypto/rand"
	"encoding/binary"
	"strings"
)

//go:generate mockgen -source=service.go -destination=mocks/service_mock.go -package=mocks
type Repository interface {
	StoreURLSInDB(originalURL, shortURL string) error
	GetShortURLFromDB(originalURL string) (string, error)
	GetOriginalURLFromDB(shortURL string) (string, error)
}
type Encoder interface {
	CryptoBase62Encode() string
}

type Services struct {
	repository Repository
	encoder    Encoder
	baseURL    string
}

func NewServices(repository Repository, encoder Encoder, baseURL string) *Services {
	return &Services{
		repository: repository,
		encoder:    encoder,
		baseURL:    baseURL,
	}
}

// GetShortURL returns the short URL
func (s Services) GetShortURL(url string) (string, error) {
	shortURL, err := s.repository.GetShortURLFromDB(url)
	if err != nil {
		shortURL = s.encoder.CryptoBase62Encode()
		err = s.repository.StoreURLSInDB(url, shortURL)
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
	_, err := rand.Read(b)
	if err != nil {
		err.Error()
	}
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
