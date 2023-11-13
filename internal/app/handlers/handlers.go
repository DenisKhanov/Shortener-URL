package handlers

import (
	"context"
	"errors"
	"github.com/DenisKhanov/shorterURL/internal/app/models"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"time"
)

// Service defines the interface for URL shortening and retrieval operations.
// It abstracts the logic for shortening URLs, fetching original URLs from shortened ones, and handling batch operations for URL shortening.
//
//go:generate mockgen -source=handlers.go -destination=mocks/handlers_mock.go -package=mocks
type Service interface {
	// GetShortURL takes original URL and returns its shortened version.
	// If the URL has already been shortened, it returns the existing shortened URL.
	// If the URL is new, it generates a new shortened URL.
	// Returns an error if the URL cannot be shortened or if any internal error occurs.
	GetShortURL(url string) (string, error)
	// GetOriginalURL takes a shortened URL and returns the original URL it points to.
	// If the shortened URL does not exist or is invalid, an error is returned.
	// Useful for redirecting shortened URLs to their original destinations.
	GetOriginalURL(shortURL string) (string, error)
	// GetBatchJSONShortURL takes a slice of models.URLRequest objects, each containing a URL to be shortened,
	// and returns a slice of models.URLResponse objects, each containing the original and shortened URL.
	// This method is intended for processing multiple URLs at once, improving efficiency for bulk operations.
	// Returns an error if any of the URLs cannot be processed or if an internal error occurs.
	GetBatchJSONShortURL(batchURLRequests []models.URLRequest) ([]models.URLResponse, error)
}

type Handlers struct {
	service Service
	DB      *pgxpool.Pool
}
type URLProcessing struct {
	URL string `json:"url"`
}
type URLProcessingResult struct {
	Result string `json:"result"`
}
type responseData struct {
	status int
	size   int
}
type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

var typeArray = [2]string{"application/json", "text/html"}

func NewHandlers(service Service, DB *pgxpool.Pool) *Handlers {
	return &Handlers{
		service: service,
		DB:      DB,
	}
}

func (h Handlers) GetShortURL(c *gin.Context) {
	linc, err := c.GetRawData()
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	lincString := string(linc)
	parsedLinc, err := url.Parse(lincString)
	if err != nil || parsedLinc.Scheme == "" || parsedLinc.Host == "" {
		c.Status(http.StatusBadRequest)
		return
	}
	shortURL, err := h.service.GetShortURL(lincString)
	if err != nil {
		if errors.Is(err, models.ErrURLFound) {
			c.String(http.StatusConflict, shortURL)
			return
		}
		c.Status(http.StatusBadRequest)
		return
	}
	c.String(http.StatusCreated, shortURL)

}
func (h Handlers) GetOriginalURL(c *gin.Context) {
	shortURL := c.Param("id")
	originURL, err := h.service.GetOriginalURL(shortURL)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.Header("Location", originURL)
	c.Status(http.StatusTemporaryRedirect)
}
func (h Handlers) GetJSONShortURL(c *gin.Context) {
	var dataURL URLProcessing
	if err := c.ShouldBindJSON(&dataURL); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	result, err := h.service.GetShortURL(dataURL.URL)
	if err != nil {
		if errors.Is(err, models.ErrURLFound) {
			c.JSON(http.StatusConflict, gin.H{"result": result})
			return
		}
		c.Status(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"result": result})
}
func (h Handlers) GetBatchJSONShortURL(c *gin.Context) {
	var batchURLRequests []models.URLRequest
	if err := c.ShouldBindJSON(&batchURLRequests); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	batchURLResponses, err := h.service.GetBatchJSONShortURL(batchURLRequests)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusCreated, batchURLResponses)
}
func (h Handlers) PingDB(c *gin.Context) {
	if h.DB != nil {
		if err := h.DB.Ping(context.Background()); err != nil {
			logrus.Error(err)
			return
		}
		c.Status(http.StatusOK)
		return
	}
	c.Status(http.StatusInternalServerError)

}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func (h Handlers) MiddlewareLogging() gin.HandlerFunc {
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
		}).Info("Обработан запрос")
	}
}
