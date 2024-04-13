// Package url provides HTTP request handlers and middleware for the URL shortening application.
package url

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

// DelUserURLs marks specified URLs as deleted asynchronously.
// Expects a JSON array of URL IDs to delete in the request body.
// Acknowledges the deletion request with HTTP status 202 Accepted.
// Returns HTTP status 400 Bad Request for malformed JSON input.
func (h *Handlers) DelUserURLs(c *gin.Context) {
	ctx := c.Request.Context()
	var URLSToDel []string
	if err := c.ShouldBindJSON(&URLSToDel); err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.AsyncDeleteUserURLs(ctx, URLSToDel); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusAccepted)
}
