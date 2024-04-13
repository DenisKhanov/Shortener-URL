// Package interceptors provides gRPC interceptors for handling authentication and authorization.
// It includes functionality to enforce authentication for specific gRPC methods.
package interceptors

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net"
)

// grpcPath defines the path for gRPC methods.
const grpcPath = "/shortener_v1.Shortener_v1/"

// subnetMethods specifies the gRPC methods that require check ip in trusted subnets.
var subnetMethods = map[string]struct{}{
	grpcPath + "GetServiceStats": {},
}

// UnaryTrustedSubnetsInterceptor creates a gRPC interceptor that checks whether the user's IP address belongs
// to trusted subnets. This interceptor is intended for use in gRPC methods that require access only from
// specific networks.
//
// Parameters:
//   - subnets: A list of trusted subnets in the format []*net.IPNet.
func UnaryTrustedSubnetsInterceptor(subnets []*net.IPNet) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{},
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if _, exist := subnetMethods[info.FullMethod]; !exist {
			return handler(ctx, req)
		}
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			values := md.Get("x-real-ip")
			if len(values) > 0 {
				ip := values[0]
				if isIPInSubnet(ip, subnets) {
					return handler(ctx, req)
				}
			}
		}
		return nil, status.Error(codes.Aborted, `this ip is not in trusted Subnet`)
	}
}

// IsIPInSubnet checks whether the specified IP address is included in the CIDR subnet.
func isIPInSubnet(ip string, subnets []*net.IPNet) bool {
	ipNet := net.ParseIP(ip)
	for _, subnet := range subnets {
		if subnet.Contains(ipNet) {
			return true
		}
	}
	return false
}
