// Package url provides the business logic for managing shortened URLs.
// It includes functionality to generate, store, retrieve, and delete URLs.
package url

import (
	"context"
	"github.com/DenisKhanov/shorterURL/internal/models"
	"github.com/sirupsen/logrus"
)

// GetBatchShortURL takes a slice of models.URLRequest objects, each containing a URL to be shortened,
// and returns a slice of models.URLResponse objects, each containing the original and shortened URL.
// This method is intended for processing multiple URLs at once, improving efficiency for bulk operations.
// Returns an error if any of the URLs cannot be processed or if an internal error occurs.
func (s ShortURLServices) GetBatchShortURL(ctx context.Context,
	batchURLRequests []models.URLRequest) ([]models.URLResponse, error) {
	shortsURL, err := s.repository.GetShortBatchURL(ctx, batchURLRequests)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	var batchURLtoStores = make(map[string]string, len(batchURLRequests))
	var batchURLResponses []models.URLResponse
	for _, value := range batchURLRequests {
		if shortURL, ok := shortsURL[value.OriginalURL]; ok {
			batchURLResponses = append(batchURLResponses, models.URLResponse{CorrelationID: value.CorrelationID, ShortURL: s.finalURLBuilder(shortURL)})
		} else {
			shortURL = s.encoder.CryptoBase62Encode()
			batchURLtoStores[shortURL] = value.OriginalURL
			batchURLResponses = append(batchURLResponses, models.URLResponse{CorrelationID: value.CorrelationID, ShortURL: s.finalURLBuilder(shortURL)})
		}
	}

	err = s.repository.StoreBatchURL(ctx, batchURLtoStores)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return batchURLResponses, nil
}
