// Package url provides the business logic for managing shortened URLs.
// It includes functionality to generate, store, retrieve, and delete URLs.
package url

import (
	"context"
	"github.com/DenisKhanov/shorterURL/internal/models"
	"github.com/sirupsen/logrus"
)

// GetUserURLs takes a slice of models.URL objects for a specific user
func (s ShortURLServices) GetUserURLs(ctx context.Context) ([]models.URL, error) {
	userURLS, err := s.repository.GetUserURLs(ctx)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	allUserShortURLs := make([]models.URL, len(userURLS))
	for i, v := range userURLS {
		allUserShortURLs[i].ShortURL = s.finalURLBuilder(v.ShortURL)
		allUserShortURLs[i].OriginalURL = v.OriginalURL
	}
	return allUserShortURLs, nil
}
