// Package interceptors provides gRPC interceptors for handling authentication and authorization.
// It includes functionality to enforce authentication for specific gRPC methods.
package interceptors

import (
	"context"
	"github.com/DenisKhanov/shorterURL/internal/auth"
	"github.com/DenisKhanov/shorterURL/internal/models"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// UnaryPublicAuthInterceptor is a gRPC interceptor that performs public authentication.
// It checks if the incoming context contains a valid token for public access.
// If the token is missing or invalid, it generates a new token and sends it in the response header.
// If the token is valid, it extracts the user ID from the token and adds it to the context.
func UnaryPublicAuthInterceptor(ctx context.Context, req interface{},
	_ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	var tokenString string
	var err error
	var userID uuid.UUID
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		values := md.Get("token")
		if len(values) > 0 {
			tokenString = values[0]
		}
	}
	if !ok || len(tokenString) == 0 || !auth.IsValidToken(tokenString) {
		logrus.Info("Token in metadata missing or isn't valid")
		tokenString, err = auth.BuildJWTString()
		if err != nil {
			logrus.Errorf("error generating token: %v", err)
			return nil, status.Errorf(codes.Unauthenticated, `error generating token: %v`, err)
		}

		if err = grpc.SendHeader(ctx, metadata.New(map[string]string{
			string(models.TokenKey): tokenString})); err != nil {
			return nil, status.Errorf(codes.Unknown, `error send token in metadata: %v`, err)
		}
		ctx = context.WithValue(ctx, models.TokenKey, tokenString)
	}

	userID, err = auth.GetUserID(tokenString)
	if err != nil {
		logrus.Error(err)
		return nil, status.Errorf(codes.Unauthenticated, `error get UUID from token: %v`, err)
	}
	ctx = context.WithValue(ctx, models.UserIDKey, userID)
	return handler(ctx, req)
}
