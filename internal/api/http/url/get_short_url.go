// Package url provides HTTP request handlers and middleware for the URL shortening application.
package url

import (
	"errors"
	"github.com/DenisKhanov/shorterURL/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
)

// GetShortURL converts a long URL to its shortened version.
// It reads the raw URL from the request body.
// Returns the shortened URL on success with HTTP status 201 Created.
// On failure, returns HTTP status 400 Bad Request for invalid input or URL format,
// or HTTP status 409 Conflict if the URL is already shortened.
func (h *Handlers) GetShortURL(c *gin.Context) {
	ctx := c.Request.Context()
	link, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	linkString := string(link)
	parsedLinc, err := url.Parse(linkString)
	if err != nil || parsedLinc.Scheme == "" || parsedLinc.Host == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("URL format isn't correct").Error()})
		return
	}
	shortURL, err := h.service.GetShortURL(ctx, linkString)
	if err != nil {
		if errors.Is(err, models.ErrURLFound) {
			c.String(http.StatusConflict, shortURL)
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.String(http.StatusCreated, shortURL)

}
