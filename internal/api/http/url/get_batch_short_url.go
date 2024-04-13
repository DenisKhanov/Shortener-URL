// Package url provides HTTP request handlers and middleware for the URL shortening application.
package url

import (
	"github.com/DenisKhanov/shorterURL/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
)

// GetBatchShortURL converts multiple URLs to their shortened versions in batch.
// Expects a JSON array of URL objects in the request body.
// Returns a JSON array of objects containing original and shortened URLs.
// On failure, sends HTTP status 500 Internal Server Error.
func (h *Handlers) GetBatchShortURL(c *gin.Context) {
	ctx := c.Request.Context()
	var batchURLRequests []models.URLRequest
	if err := c.ShouldBindJSON(&batchURLRequests); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for _, req := range batchURLRequests {
		parsedLinc, err := url.Parse(req.OriginalURL)
		if err != nil || parsedLinc.Scheme == "" || parsedLinc.Host == "" {
			c.String(http.StatusBadRequest, "URL format isn't correct")
			return
		}
	}
	batchURLResponses, err := h.service.GetBatchShortURL(ctx, batchURLRequests)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, batchURLResponses)
}
