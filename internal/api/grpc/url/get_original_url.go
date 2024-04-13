// Package url package provides functionality for handling gRPC communication related to URL shortening.
// It includes interfaces and structs for defining gRPC services and servers, as well as methods
// for interacting with the URL shortening service.
package url

import (
	"context"
	"errors"
	"github.com/DenisKhanov/shorterURL/internal/models"
	proto "github.com/DenisKhanov/shorterURL/pkg/shortener_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

// GetOriginalURL method in the ShortenerServer struct handles gRPC requests to
// retrieve the original URL corresponding to a shortened URL. It extracts the
// short URL from the request, then invokes the GetOriginalURL method of the
// service layer to retrieve the corresponding original URL. If the retrieval
// is successful, it constructs a response containing the original URL and returns
// it along with a status error with the OK code and a message indicating that the
// original URL was successfully obtained.
//
// If an error occurs during the retrieval process, it checks if the error indicates
// that the URL has been deleted. If so, it returns a status error with the NotFound
// code and an appropriate error message. Otherwise, it returns a status error with
// the InvalidArgument code and an appropriate error message.
func (s *ShortenerServer) GetOriginalURL(ctx context.Context,
	in *proto.GetOriginalURLRequest) (*proto.GetOriginalURLResponse, error) {
	var response proto.GetOriginalURLResponse
	parts := strings.Split(in.ShortUrl, "/")
	shortURL := parts[len(parts)-1]
	originURL, err := s.service.GetOriginalURL(ctx, shortURL)
	if err != nil {
		if errors.Is(err, models.ErrURLDeleted) {

			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		return nil, status.Errorf(codes.InvalidArgument, `error: %v`, err)
	}
	response.OriginalUrl = originURL
	return &response, status.Error(codes.OK, `original url`)
}
