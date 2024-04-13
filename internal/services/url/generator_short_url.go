// Package url provides the business logic for managing shortened URLs.
// It includes functionality to generate, store, retrieve, and delete URLs.
package url

import (
	"crypto/rand"
	"encoding/binary"
	"github.com/sirupsen/logrus"
	"net/url"
)

// finalURLBuilder the function combines the base url and the shortened url into a single link
func (s ShortURLServices) finalURLBuilder(shortURL string) string {
	resultURL, err := url.JoinPath(s.baseURL, shortURL)
	if err != nil {
		logrus.Error(err)
	}
	return resultURL
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
