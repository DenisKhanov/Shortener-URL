// Package interceptors provides gRPC interceptors for handling authentication and authorization.
// It includes functionality to enforce authentication for specific gRPC methods.
package interceptors

import (
	"context"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"time"
)

// UnaryLoggerInterceptor is a gRPC interceptor that logs information about unary RPC requests and responses.
// It measures the duration of the request handling process and logs details such as the method, status code, error message (if any), and duration.
func UnaryLoggerInterceptor(ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	start := time.Now()

	resp, err := handler(ctx, req)

	duration := time.Since(start)

	inf := status.Convert(err)
	var statusCode, errMsg string

	if err != nil {
		statusCode = inf.Code().String()
		errMsg = inf.Message()
	} else {
		statusCode = "OK"
	}

	//logging request/response info
	logrus.WithFields(logrus.Fields{
		"method":   info.FullMethod,
		"status":   statusCode,
		"errMsg":   errMsg,
		"duration": duration,
	}).Info("Обработан запрос gRPC")
	return resp, err
}
