// Package url package provides functionality for handling gRPC communication related to URL shortening.
// It includes interfaces and structs for defining gRPC services and servers, as well as methods
// for interacting with the URL shortening service.
package url

import (
	"context"
	"github.com/DenisKhanov/shorterURL/internal/models"
	proto "github.com/DenisKhanov/shorterURL/pkg/shortener_v1"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/url"
)

// GetBatchShortURL method in the ShortenerServer struct handles gRPC requests to
// compress multiple URLs into shortened versions. It parses each original URL from
// the request, validates its format, and constructs a batch of URL requests. It then
// invokes the GetBatchShortURL method of the service layer to compress the batch of URLs.
// If successful, it constructs a response containing the compressed URLs and returns it
// along with a status error with the OK code and a message indicating that all URLs have
// been compressed.
//
// If an error occurs during the compression process, it logs the error, constructs an
// appropriate error message, and returns a status error with the Internal code.
func (s *ShortenerServer) GetBatchShortURL(ctx context.Context,
	in *proto.GetBatchShortURLRequest) (*proto.GetBatchShortURLResponse, error) {
	var response proto.GetBatchShortURLResponse
	batchURLRequests := make([]models.URLRequest, len(in.BatchUrlRequests))
	for i, req := range in.BatchUrlRequests {
		parsedLinc, err := url.Parse(req.OriginalUrl)
		if err != nil || parsedLinc.Scheme == "" || parsedLinc.Host == "" {
			return nil, status.Error(codes.InvalidArgument, `URL format isn't correct`)
		}
		batchURLRequests[i] = models.URLRequest{
			CorrelationID: req.CorrelationId,
			OriginalURL:   req.OriginalUrl,
		}
	}
	batchURLResponses, err := s.service.GetBatchShortURL(ctx, batchURLRequests)
	if err != nil {
		logrus.Error(err)
		return nil, status.Errorf(codes.Internal, `error: %v`, err)
	}
	resultBatchURLResponses := make([]*proto.URLResponse, len(in.BatchUrlRequests))
	for i, res := range batchURLResponses {
		resultBatchURLResponses[i] = &proto.URLResponse{
			CorrelationId: res.CorrelationID,
			ShortUrl:      res.ShortURL,
		}
	}
	response.BatchUrlResponses = resultBatchURLResponses
	return &response, status.Error(codes.OK, `all URLs have bin compressed`)
}
