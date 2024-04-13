// Package url package provides functionality for handling gRPC communication related to URL shortening.
// It includes interfaces and structs for defining gRPC services and servers, as well as methods
// for interacting with the URL shortening service.
package url

import (
	"context"
	proto "github.com/DenisKhanov/shorterURL/pkg/shortener_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetServiceStats method within the ShortenerServer struct handles gRPC
// requests to retrieve statistics about the service. It delegates the retrieval
// process to the service layer's GetServiceStats method. If the retrieval is successful,
// it constructs a response containing the retrieved statistics and returns it along
// with a status error with the OK code and a message indicating that the statistics
// were successfully obtained.
//
// If an error occurs during the retrieval process, it returns a status error with
// the Internal code and an appropriate error message.
func (s *ShortenerServer) GetServiceStats(ctx context.Context,
	in *proto.GetServiceStatsRequest) (*proto.GetServiceStatsResponse, error) {
	var response proto.GetServiceStatsResponse
	stats, err := s.service.GetServiceStats(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, `error: %v`, err)
	}
	response.Stats = &proto.Stats{
		CountUrls:  stats.CountURLs,
		CountUsers: stats.CountUsers,
	}
	return &response, status.Error(codes.OK, `stats got`)
}
