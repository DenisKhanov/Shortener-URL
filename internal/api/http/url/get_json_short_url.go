// Package url provides HTTP request handlers and middleware for the URL shortening application.
package url

import (
	"errors"
	"github.com/DenisKhanov/shorterURL/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
)

// GetJSONShortURL converts a long URL to its shortened version using JSON input.
// Expects a JSON object with a 'URL' field in the request body.
// Returns a JSON object containing the shortened URL on success.
// Sends HTTP status 400 Bad Request for malformed JSON or invalid URL,
// or HTTP status 409 Conflict if the URL is already shortened.
func (h *Handlers) GetJSONShortURL(c *gin.Context) {
	ctx := c.Request.Context()

	var dataURL URLProcessing
	if err := c.ShouldBindJSON(&dataURL); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON"})
		return
	}

	parsedLinc, err := url.Parse(dataURL.URL)
	if err != nil || parsedLinc.Scheme == "" || parsedLinc.Host == "" {
		c.String(http.StatusBadRequest, "URL format isn't correct")
		return
	}

	result, err := h.service.GetShortURL(ctx, dataURL.URL)
	if err != nil {
		if errors.Is(err, models.ErrURLFound) {
			c.JSON(http.StatusConflict, gin.H{"result": result})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"result": result})
}
