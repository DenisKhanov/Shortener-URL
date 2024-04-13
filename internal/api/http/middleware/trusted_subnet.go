// Package middleware provides HTTP middleware for handlers.
package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
)

// TrustedSubnet returns a Gin middleware function that checks if the client's IP is within the trusted subnets.
//
// This middleware function retrieves the client's IP using the ClientIP method of the Gin context.
// It then checks if the IP is within the trusted subnets using the isIPInSubnet method of the Handlers struct.
// If the IP is not within the trusted subnets, it logs an error and responds with a 403 Forbidden status code along with an error message.
// Otherwise, it calls the Next function to proceed with the next middleware or route handler in the chain.
func TrustedSubnet(subnets []*net.IPNet) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !isIPInSubnet(ip, subnets) {
			logrus.Error(errors.New("this ip is not in trusted Subnet"))
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "this ip is not in trusted Subnet"})
		}
		c.Next()
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
