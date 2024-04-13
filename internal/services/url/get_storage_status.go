// Package url provides the business logic for managing shortened URLs.
// It includes functionality to generate, store, retrieve, and delete URLs.
package url

import (
	"context"
	"github.com/sirupsen/logrus"
)

// GetStorageStatus method returns an error after checking the status of the storage.
// It performs a check of the storage's availability by calling the Ping method from
// the repository. If the storage is unavailable or if an error occurs during the check,
// the method returns an error.
func (s ShortURLServices) GetStorageStatus(ctx context.Context) error {
	if err := s.repository.Ping(ctx); err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}
