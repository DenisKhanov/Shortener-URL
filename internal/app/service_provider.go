package app

import (
	url3 "github.com/DenisKhanov/shorterURL/internal/api/grpc/url"
	url4 "github.com/DenisKhanov/shorterURL/internal/api/http/url"
	url2 "github.com/DenisKhanov/shorterURL/internal/repositories/url"
	"github.com/DenisKhanov/shorterURL/internal/services/url"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

// serviceProvider manages the dependency injection for http_shortener-related components.
type serviceProvider struct {
	shortenerRepository url.Repository        // Repository for http_shortener-related data
	shortenerService    url4.Service          // Service for http_shortener-related operations
	shortenerHandler    *url4.Handlers        // Handler for http_shortener-related HTTP endpoints
	shortenerGRPC       *url3.ShortenerServer //GRPC for http_shortener-related operations
}

// newServiceProvider creates a new instance of the service provider.
func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

// ShortenerRepository returns the repository for user-related data.
// If dbPool is nil, it initializes an in-memory repository, otherwise initializes a database repository.
func (s *serviceProvider) ShortenerRepository(dbPool *pgxpool.Pool, storagePath string) url.Repository {
	var err error
	if s.shortenerRepository == nil {
		if dbPool == nil {
			s.shortenerRepository = url2.NewURLInMemoryRepo(storagePath)
		} else {
			if s.shortenerRepository, err = url2.NewURLInDBRepo(dbPool); err != nil {
				//TODO лучше вернуть ошибку из метода и обработать ее выше
				logrus.Fatal(err)
			}
		}
	}
	return s.shortenerRepository
}

// ShortenerService returns the service for user-related operations.
func (s *serviceProvider) ShortenerService(dbPool *pgxpool.Pool, baseUrl, storagePath string) url4.Service {
	if s.shortenerService == nil {
		s.shortenerService = url.NewShortURLServices(
			s.ShortenerRepository(dbPool, storagePath),
			url.ShortURLServices{},
			baseUrl,
		)
	}
	return s.shortenerService
}

// ShortenerHandler returns the handler for user-related HTTP endpoints.
func (s *serviceProvider) ShortenerHandler(dbPool *pgxpool.Pool, baseUrl, storagePath, subnetsStr string) *url4.Handlers {
	if s.shortenerHandler == nil {
		userHandler := url4.NewHandlers(s.ShortenerService(dbPool, baseUrl, storagePath), subnetsStr)
		s.shortenerHandler = userHandler
	}
	return s.shortenerHandler
}

// ShortenerGRPC returns the handler for user-related HTTP endpoints.
func (s *serviceProvider) ShortenerGRPC() *url3.ShortenerServer {
	if s.shortenerGRPC == nil {
		shortenerGRPC := url3.NewShortenerServer(s.shortenerService)
		s.shortenerGRPC = shortenerGRPC
	}
	return s.shortenerGRPC
}
