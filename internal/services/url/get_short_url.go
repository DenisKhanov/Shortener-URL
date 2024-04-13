// Package url provides the business logic for managing shortened URLs.
// It includes functionality to generate, store, retrieve, and delete URLs.
package url

import (
	"context"
	"github.com/DenisKhanov/shorterURL/internal/models"
)

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
