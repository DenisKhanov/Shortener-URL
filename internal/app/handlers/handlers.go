package handlers

import (
	"compress/gzip"
	"context"
	"errors"
	"github.com/DenisKhanov/shorterURL/internal/app/auth"
	"github.com/DenisKhanov/shorterURL/internal/app/models"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"strings"
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
	GetShortURL(ctx context.Context, url string) (string, error)
	// GetOriginalURL takes a shortened URL and returns the original URL it points to.
	// If the shortened URL does not exist or is invalid, an error is returned.
	// Useful for redirecting shortened URLs to their original destinations.
	GetOriginalURL(ctx context.Context, shortURL string) (string, error)
	// GetBatchShortURL takes a slice of models.URLRequest objects, each containing a URL to be shortened,
	// and returns a slice of models.URLResponse objects, each containing the original and shortened URL.
	// This method is intended for processing multiple URLs at once, improving efficiency for bulk operations.
	// Returns an error if any of the URLs cannot be processed or if an internal error occurs.
	GetBatchShortURL(ctx context.Context, batchURLRequests []models.URLRequest) ([]models.URLResponse, error)
	// GetUserURLS takes a slice of models.URL objects for a specific user
	GetUserURLS(ctx context.Context) ([]models.URL, error)
}

type Handlers struct {
	service Service
	DB      *pgxpool.Pool
}
type URLProcessing struct {
	URL string `json:"url"`
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
	ctx := c.Request.Context()
	link, err := c.GetRawData()
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	linkString := string(link)
	parsedLinc, err := url.Parse(linkString)
	if err != nil || parsedLinc.Scheme == "" || parsedLinc.Host == "" {
		c.Status(http.StatusBadRequest)
		return
	}
	shortURL, err := h.service.GetShortURL(ctx, linkString)
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
	ctx := c.Request.Context()
	shortURL := c.Param("id")
	originURL, err := h.service.GetOriginalURL(ctx, shortURL)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.Header("Location", originURL)
	c.Status(http.StatusTemporaryRedirect)
}
func (h Handlers) GetJSONShortURL(c *gin.Context) {
	ctx := c.Request.Context()
	var dataURL URLProcessing
	if err := c.ShouldBindJSON(&dataURL); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	result, err := h.service.GetShortURL(ctx, dataURL.URL)
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
func (h Handlers) GetBatchShortURL(c *gin.Context) {
	ctx := c.Request.Context()
	var batchURLRequests []models.URLRequest
	if err := c.ShouldBindJSON(&batchURLRequests); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	batchURLResponses, err := h.service.GetBatchShortURL(ctx, batchURLRequests)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusCreated, batchURLResponses)
}
func (h Handlers) GetUserURLS(c *gin.Context) {
	ctx := c.Request.Context()
	fullShortUserURLS, err := h.service.GetUserURLS(ctx)
	if err != nil {
		c.Status(http.StatusNoContent)
		return
	}
	c.JSON(http.StatusCreated, fullShortUserURLS)
}
func (h Handlers) PingDB(c *gin.Context) {
	ctx := c.Request.Context()
	if h.DB != nil {
		if err := h.DB.Ping(ctx); err != nil {
			logrus.Error(err)
			c.Status(http.StatusInternalServerError)
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

type compressWriter struct {
	gin.ResponseWriter
	Writer *gzip.Writer
}

func (c *compressWriter) Write(data []byte) (int, error) {
	return c.Writer.Write(data)
}
func (c *compressWriter) Close() error {
	return c.Writer.Close()
}
func (c *compressWriter) WriteString(s string) (int, error) {
	return c.Writer.Write([]byte(s))
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
func (h Handlers) MiddlewareCompress() gin.HandlerFunc {
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
func (h Handlers) MiddlewareAuthPublic() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string
		var err error
		var userID uint32

		tokenString, err = c.Cookie("user_token")
		// если токен не найден в куке, то генерируем новый и добавляем его в куки
		if err != nil || !auth.IsValidToken(tokenString) {
			logrus.Info("Cookie not found or token in cookie not found")
			tokenString, err = auth.BuildJWTString()
			if err != nil {
				logrus.Errorf("error generating token: %v", err)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
			c.SetCookie("user_token", tokenString, 0, "/", "", false, true)
		}
		userID, err = auth.GetUserID(tokenString)
		if err != nil {
			logrus.Error(err)
			return
		}

		ctx := context.WithValue(c.Request.Context(), "userID", userID)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
func (h Handlers) MiddlewareAuthPrivate() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("user_token")
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		userID, err := auth.GetUserID(tokenString)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		ctx := context.WithValue(c.Request.Context(), "userID", userID)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
