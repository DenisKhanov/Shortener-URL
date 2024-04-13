// Package url provides the business logic for managing shortened URLs.
// It includes functionality to generate, store, retrieve, and delete URLs.
package url

import (
	"context"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"time"
)

// AsyncDeleteUserURLs async runs requests to DB for mark user URLs as deleted
func (s ShortURLServices) AsyncDeleteUserURLs(ctx context.Context, URLSToDel []string) error {
	var g errgroup.Group

	g.Go(func() error {
		asyncCtx, cancel := context.WithTimeout(ctx, time.Minute)
		defer cancel()
		if err := s.repository.MarkURLsAsDeleted(asyncCtx, URLSToDel); err != nil {
			logrus.Error(err)
			return err
		}
		return nil
	})
	// Ожидание завершения работы всех горутин и проверка ошибок
	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}
