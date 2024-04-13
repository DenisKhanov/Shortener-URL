// Package interceptors provides gRPC interceptors for handling authentication and authorization.
// It includes functionality to enforce authentication for specific gRPC methods.
package interceptors

import (
	"context"
	"github.com/DenisKhanov/shorterURL/internal/auth"
	"github.com/DenisKhanov/shorterURL/internal/models"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// grpcHandlersPath defines the path for gRPC methods.
const grpcHandlersPath = "/shortener_v1.Shortener_v1/"

// authMethods specifies the gRPC methods that require authentication.
var authMethods = map[string]struct{}{
	grpcHandlersPath + "GetUserURLs": {},
	grpcHandlersPath + "DelUserURLs": {},
}

// UnaryPrivateAuthInterceptor is a gRPC interceptor that enforces authentication for specific unary RPCs.
// It checks if the incoming context contains a valid token for accessing the specified methods.
// If the token is valid, it extracts the user ID from the token and adds it to the context.
// If the token is missing or invalid, it returns an error.
func UnaryPrivateAuthInterceptor(ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if _, exist := authMethods[info.FullMethod]; !exist {
		return handler(ctx, req)
	}
	var tokenString string
	var err error
	var userID uuid.UUID
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		values := md.Get("token")
		if len(values) > 0 {
			// ключ содержит слайс строк, получаем первую строку
			tokenString = values[0]
		}
	}
	if !ok || len(tokenString) == 0 {
		return nil, status.Error(codes.InvalidArgument, `missing token`)
	}
	userID, err = auth.GetUserID(tokenString)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, `invalid token`)
	}

	ctx = context.WithValue(ctx, models.UserIDKey, userID)
	return handler(ctx, req)
}
