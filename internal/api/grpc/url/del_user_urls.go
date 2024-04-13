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

// DelUserURLs method in the ShortenerServer struct handles gRPC requests
// to asynchronously mark user URLs as deleted. It invokes the AsyncDeleteUserURLs
// method of the service layer to perform the deletion operation. If an error occurs
// during the deletion process, it constructs an appropriate error message and returns
// a status error with the Internal code. Otherwise, it returns a status error with the
// OK code and a message indicating that the URLs have been marked as deleted.
func (s *ShortenerServer) DelUserURLs(ctx context.Context,
	in *proto.DelUserURLsRequest) (*proto.DelUserURLsResponse, error) {
	if err := s.service.AsyncDeleteUserURLs(ctx, in.UrlsToDel); err != nil {
		return nil, status.Errorf(codes.Internal, `error: %v`, err)
	}
	return nil, status.Error(codes.OK, `URLs marked as deleted`)
}
