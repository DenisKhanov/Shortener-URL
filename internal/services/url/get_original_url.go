// Package url provides the business logic for managing shortened URLs.
// It includes functionality to generate, store, retrieve, and delete URLs.
package url

import "context"

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
