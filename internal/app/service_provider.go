package app

import (
	"github.com/DenisKhanov/shorterURL/internal/handlers"
	"github.com/DenisKhanov/shorterURL/internal/repositories"
	"github.com/DenisKhanov/shorterURL/internal/services"
	"github.com/jackc/pgx/v5/pgxpool"
)

// serviceProvider manages the dependency injection for user-related components.
type serviceProvider struct {
	userRepository services.Repository // Repository for user-related data
	userService    handlers.Service    // Service for user-related operations
	userHandler    *handlers.Handlers  // Handler for user-related HTTP endpoints
}

// newServiceProvider creates a new instance of the service provider.
func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

// UserRepository returns the repository for user-related data.
// If dbPool is nil, it initializes an in-memory repository, otherwise initializes a database repository.
func (s *serviceProvider) UserRepository(dbPool *pgxpool.Pool, storagePath string) services.Repository {
	if s.userRepository == nil {
		if dbPool == nil {
			s.userRepository = repositories.NewURLInMemoryRepo(storagePath)
		} else {
			s.userRepository = repositories.NewURLInDBRepo(dbPool)
		}
	}
	return s.userRepository
}

// UserService returns the service for user-related operations.
func (s *serviceProvider) UserService(dbPool *pgxpool.Pool, baseUrl, storagePath string) handlers.Service {
	if s.userService == nil {
		s.userService = services.NewShortURLServices(
			s.UserRepository(dbPool, storagePath),
			services.ShortURLServices{},
			baseUrl,
		)
	}
	return s.userService
}

// UserHandler returns the handler for user-related HTTP endpoints.
func (s *serviceProvider) UserHandler(dbPool *pgxpool.Pool, baseUrl, storagePath, subnetsStr string) (*handlers.Handlers, error) {
	if s.userHandler == nil {
		userHandler, err := handlers.NewHandlers(s.UserService(dbPool, baseUrl, storagePath), dbPool, subnetsStr)
		if err != nil {
			return nil, err
		}
		s.userHandler = userHandler
	}
	return s.userHandler, nil
}
