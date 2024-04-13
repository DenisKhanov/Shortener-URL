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

//TODO добавить вывод статуса удаления URL

// GetUserURLs method within the ShortenerServer struct handles requests in gRPC format to
// retrieve user-specific URLs. It delegates the retrieval process to the service layer's GetUserURLs
// method. If no URLs are found for the user, it returns a status error with the NotFound code and an
// appropriate error message. Otherwise, it constructs a response containing all user URLs in
// the appropriate format and returns it along with a status error with the OK code and a message
// indicating that all user URLs have been successfully retrieved.
func (s *ShortenerServer) GetUserURLs(ctx context.Context,
	in *proto.GetUserURLsRequest) (*proto.GetUserURLsResponse, error) {
	var response proto.GetUserURLsResponse
	allUserShortURLs, err := s.service.GetUserURLs(ctx)
	if err != nil {
		return nil, status.Error(codes.NotFound, `your short URLs not found`)
	}
	resultAllUserShortURLs := make([]*proto.URL, len(allUserShortURLs))
	for i, res := range allUserShortURLs {
		resultAllUserShortURLs[i] = &proto.URL{
			ShortUrl:    res.ShortURL,
			OriginalUrl: res.OriginalURL,
		}
	}
	response.UserUrls = resultAllUserShortURLs
	return &response, status.Error(codes.OK, `your all short URLs`)
}
