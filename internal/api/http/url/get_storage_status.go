// Package url provides HTTP request handlers and middleware for the URL shortening application.
package url

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

// GetStorageStatus method in the Handlers struct handles the HTTP request to check the status of the storage backend.
// It invokes the GetStorageStatus method of the service layer to perform the check. If an error occurs during the check,
// it logs the error and returns a JSON response with an internal server error status code along with the error message.
// Otherwise, it responds with a status code of 200 (OK) to indicate that the storage is accessible.
func (h *Handlers) GetStorageStatus(c *gin.Context) {
	ctx := c.Request.Context()
	if err := h.service.GetStorageStatus(ctx); err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
