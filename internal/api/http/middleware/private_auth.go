// Package middleware provides HTTP middleware for handlers.
package middleware

import (
	"context"
	"github.com/DenisKhanov/shorterURL/internal/auth"
	"github.com/DenisKhanov/shorterURL/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

// AuthPrivate provides authentication middleware for private routes.
// It checks the user token and only allows access if the token is valid.
// This middleware ensures that only authenticated users can access certain routes.
func AuthPrivate() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string
		var err error
		var userID uuid.UUID
		tokenString, err = c.Cookie("user_token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		}
		userID, err = auth.GetUserID(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		}
		ctx := context.WithValue(c.Request.Context(), models.UserIDKey, userID)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
