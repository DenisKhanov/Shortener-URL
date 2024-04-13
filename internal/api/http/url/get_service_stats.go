// Package url provides HTTP request handlers and middleware for the URL shortening application.
package url

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetServiceStats first retrieves the statistics using the GetServiceStats method of the service.
// If an error occurs during the retrieval process, it responds with a 403 Forbidden status code and an error message.
// Otherwise, it responds with a 200 OK status code and the retrieved statistics in JSON format.
func (h *Handlers) GetServiceStats(c *gin.Context) {
	fmt.Println("Я тут")
	ctx := c.Request.Context()
	stats, err := h.service.GetServiceStats(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}
