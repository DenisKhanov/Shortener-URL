// Package url provides HTTP request handlers and middleware for the URL shortening application.
package url

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

//TODO добавить вывод статуса удаления URL

// GetUserURLS retrieves all URLs associated with the current user.
// Does not require any parameters; user identification is from the context.
// Returns a JSON array of URLs associated with the user.
// If no URLs are found, returns HTTP status 204 No Content.
func (h *Handlers) GetUserURLS(c *gin.Context) {
	ctx := c.Request.Context()
	allUserShortURLs, err := h.service.GetUserURLs(ctx)
	if err != nil {
		c.JSON(http.StatusNoContent, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, allUserShortURLs)
}
