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
	"net/url"
)

// GetShortURL method within the ShortenerServer struct handles gRPC requests
// to retrieve a shortened URL corresponding to a given original URL. It first
// validates the format of the original URL provided in the request. If the URL
// format is incorrect, it returns a status error with the InvalidArgument code
// and an appropriate error message.
//
// If the URL format is valid, it delegates the retrieval process to the service
// layer's GetShortURL method. If a shortened URL is found in the database, it
// constructs a response containing the shortened URL and returns it along with
// a status error with the OK code and a message indicating that the shortened URL was found.
//
// If no shortened URL is found in the database, it constructs a response containing
// the newly generated shortened URL and returns it along with a status error with the
// OK code and a message indicating that the request was completed successfully.
// If an unexpected error occurs during the process, it returns a status error with
// the Unknown code and an appropriate error message.
func (s *ShortenerServer) GetShortURL(ctx context.Context,
	in *proto.GetShortURLRequest) (*proto.GetShortURLResponse, error) {
	var response proto.GetShortURLResponse
	linkString := in.OriginalUrl
	parsedLinc, err := url.Parse(linkString)
	if err != nil || parsedLinc.Scheme == "" || parsedLinc.Host == "" {
		return nil, status.Error(codes.InvalidArgument, `URL format isn't correct`)
	}
	shortURL, err := s.service.GetShortURL(ctx, linkString)
	if err != nil {
		if errors.Is(err, models.ErrURLFound) {
			response.ShortUrl = shortURL
			return &response, status.Error(codes.OK, `short URL found in database`)
		}
		return nil, status.Errorf(codes.Unknown, `error: %v`, err)
	}
	response.ShortUrl = shortURL
	return &response, status.Error(codes.OK, `request completed`)
}
