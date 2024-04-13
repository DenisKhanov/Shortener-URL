// Package middleware provides HTTP middleware for handlers.
package middleware

import (
	"compress/gzip"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// compressWriter is a custom response writer that performs gzip compression.
type compressWriter struct {
	gin.ResponseWriter
	Writer *gzip.Writer
}

// compressWriter provides a custom response writer for gzip compression.
func (c *compressWriter) Write(data []byte) (int, error) {
	return c.Writer.Write(data)
}

// Close overrides the Close method to close the gzip writer.
func (c *compressWriter) Close() error {
	return c.Writer.Close()
}

// WriteString overrides the WriteString method to write a string to the gzip writer.
func (c *compressWriter) WriteString(s string) (int, error) {
	return c.Writer.Write([]byte(s))
}

// GZIPCompress provides a compression middleware using gzip.
// It checks the 'Accept-Encoding' header of incoming requests and applies gzip compression if applicable.
// This middleware optimizes response size and speed, improving overall performance.
func GZIPCompress() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.Contains(c.GetHeader("Accept-Encoding"), "gzip") {
			gz := gzip.NewWriter(c.Writer)
			defer gz.Close()
			c.Writer = &compressWriter{Writer: gz, ResponseWriter: c.Writer}
			c.Header("Content-Encoding", "gzip")
		}
		// Проверяем, сжат ли запрос
		if strings.Contains(c.GetHeader("Content-Encoding"), "gzip") {
			reader, err := gzip.NewReader(c.Request.Body)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid gzip body"})
				return
			}
			defer reader.Close()
			c.Request.Body = reader
		}
		c.Next()
	}
}
