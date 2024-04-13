// Package url provides HTTP request handlers and middleware for the URL shortening application.
package url

import (
	"errors"
	"github.com/DenisKhanov/shorterURL/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetOriginalURL retrieves the original URL from a shortened URL ID.
// The shortened URL ID is expected as a URL parameter.
// Redirects to the original URL using HTTP 307 Temporary Redirect.
// Returns HTTP status 410 Gone if the URL is marked as deleted,
// or HTTP status 400 Bad Request for other errors.
func (h *Handlers) GetOriginalURL(c *gin.Context) {
	ctx := c.Request.Context()
	shortURL := c.Param("id")
	originURL, err := h.service.GetOriginalURL(ctx, shortURL)
	if err != nil {
		if errors.Is(err, models.ErrURLDeleted) {
			c.JSON(http.StatusGone, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Header("Location", originURL)
	c.Status(http.StatusTemporaryRedirect)
}
