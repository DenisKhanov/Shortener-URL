// Package url provides the business logic for managing shortened URLs.
// It includes functionality to generate, store, retrieve, and delete URLs.
package url

import (
	"context"
	"github.com/DenisKhanov/shorterURL/internal/models"
)

// GetServiceStats retrieves the statistics of URLs and users from the service's repository.
//
// This method delegates the retrieval of statistics to the repository's GetStats method.
// It then returns the obtained statistics and any error encountered during the retrieval process.
func (s ShortURLServices) GetServiceStats(ctx context.Context) (models.Stats, error) {
	stats, err := s.repository.GetStats(ctx)
	if err != nil {
		return models.Stats{}, err
	}
	return stats, err
}
