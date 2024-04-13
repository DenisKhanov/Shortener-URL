// Package url package provides functionality for handling gRPC communication related to URL shortening.
// It includes interfaces and structs for defining gRPC services and servers, as well as methods
// for interacting with the URL shortening service.
package url

import (
	"context"
	proto "github.com/DenisKhanov/shorterURL/pkg/shortener_v1"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetStorageStatus method in the ShortenerServer struct handles the gRPC request
// to check the status of the storage backend. It invokes the GetStorageStatus method
// of the service layer to perform the check. If an error occurs during the check,
// it logs the error and returns an error status with the Unavailable code along
// with an error message. Otherwise, it returns an error status with the OK code along
// with a message indicating that the database connection is enabled.
func (s *ShortenerServer) GetStorageStatus(ctx context.Context,
	_ *proto.GetStorageStatusRequest) (*proto.GetStorageStatusResponse, error) {
	if err := s.service.GetStorageStatus(ctx); err != nil {
		logrus.Error(err)
		return nil, status.Error(codes.Unavailable, `error DB connection`)
	}
	return nil, status.Error(codes.OK, `DB connection is enable`)
}
