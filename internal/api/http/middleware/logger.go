// Package middleware provides HTTP middleware for handlers.
package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"time"
)

// LogrusLog provides a logging middleware for Gin.
// It logs details about each request including the URL, method, response status, duration, and size.
// This middleware is useful for monitoring and debugging purposes.
func LogrusLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Запуск таймера
		start := time.Now()

		// Обработка запроса
		c.Next()

		// Измерение времени обработки
		duration := time.Since(start)

		// Получение статуса ответа и размера
		status := c.Writer.Status()
		size := c.Writer.Size()

		// Логирование информации о запросе
		logrus.WithFields(logrus.Fields{
			"url":      c.Request.URL.RequestURI(),
			"method":   c.Request.Method,
			"status":   status,
			"duration": duration,
			"size":     size,
		}).Info("Обработан запрос HTTP")
	}
}
