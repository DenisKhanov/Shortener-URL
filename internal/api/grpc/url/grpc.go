// Package url package provides functionality for handling gRPC communication related to URL shortening.
// It includes interfaces and structs for defining gRPC services and servers, as well as methods
// for interacting with the URL shortening service.
package url

import (
	"context"
	"github.com/DenisKhanov/shorterURL/internal/models"
	url2 "github.com/DenisKhanov/shorterURL/internal/services/url"
	proto "github.com/DenisKhanov/shorterURL/pkg/shortener_v1"
)

//go:generate mockgen -source=grpc.go -destination=mocks/grpc_mock.go -package=mocks
type Service interface {
	// GetStorageStatus checks the database connection or repository created.
	GetStorageStatus(ctx context.Context) error
	// GetShortURL takes original URL and returns its shortened version.
	// If the URL has already been shortened, it returns the existing shortened URL.
	// If the URL is new, it generates a new shortened URL.
	// Returns an error if the URL cannot be shortened or if any internal error occurs.
	GetShortURL(ctx context.Context, url string) (string, error)
	// GetOriginalURL takes a shortened URL and returns the original URL it points to.
	// If the shortened URL does not exist or is invalid, an error is returned.
	// Useful for redirecting shortened URLs to their original destinations.
	GetOriginalURL(ctx context.Context, shortURL string) (string, error)
	// GetBatchShortURL takes a slice of models.URLRequest objects, each containing a URL to be shortened,
	// and returns a slice of models.URLResponse objects, each containing the original and shortened URL.
	// This method is intended for processing multiple URLs at once, improving efficiency for bulk operations.
	// Returns an error if any of the URLs cannot be processed or if an internal error occurs.
	GetBatchShortURL(ctx context.Context, batchURLRequests []models.URLRequest) ([]models.URLResponse, error)
	// GetUserURLs takes a slice of models.URL objects for a specific user
	GetUserURLs(ctx context.Context) ([]models.URL, error)
	// AsyncDeleteUserURLs async runs requests to DB for mark user URLs as deleted
	AsyncDeleteUserURLs(ctx context.Context, URLSToDel []string) error
	// GetServiceStats retrieves the statistics of URLs and users from the service's repository.
	GetServiceStats(ctx context.Context) (models.Stats, error)
}

// checking interface compliance at the compiler level
var _ Service = (*url2.ShortURLServices)(nil)

// ShortenerServer поддерживает все необходимые методы сервера.
type ShortenerServer struct {
	// нужно встраивать тип pb.Unimplemented<TypeName>
	// для совместимости с будущими версиями
	proto.UnimplementedShortenerV1Server
	service Service
}

// NewShortenerServer function creates a new instance of the ShortenerServer struct with the
// provided service. It initializes the service field of the ShortenerServer struct with the given
// service instance and returns a pointer to the newly created ShortenerServer instance.
func NewShortenerServer(service Service) *ShortenerServer {
	return &ShortenerServer{
		service: service,
	}
}
