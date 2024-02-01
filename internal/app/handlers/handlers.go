package handlers

import (
	"compress/gzip"
	"context"
	"errors"
	"github.com/DenisKhanov/shorterURL/internal/app/auth"
	"github.com/DenisKhanov/shorterURL/internal/app/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	// AsyncDeleteUserURLs async runs requests to DB for mark user URLs as deleted
	AsyncDeleteUserURLs(ctx context.Context, URLSToDel []string)
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

func NewHandlers(service Service, DB *pgxpool.Pool) *Handlers {
	return &Handlers{
		service: service,
		DB:      DB,
	}
}

// GetShortURL converts a long URL to its shortened version.
// It reads the raw URL from the request body.
// Returns the shortened URL on success with HTTP status 201 Created.
// On failure, returns HTTP status 400 Bad Request for invalid input or URL format,
// or HTTP status 409 Conflict if the URL is already shortened.
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

// GetOriginalURL retrieves the original URL from a shortened URL ID.
// The shortened URL ID is expected as a URL parameter.
// Redirects to the original URL using HTTP 307 Temporary Redirect.
// Returns HTTP status 410 Gone if the URL is marked as deleted,
// or HTTP status 400 Bad Request for other errors.
func (h Handlers) GetOriginalURL(c *gin.Context) {
	ctx := c.Request.Context()
	shortURL := c.Param("id")
	originURL, err := h.service.GetOriginalURL(ctx, shortURL)
	if err != nil {
		if errors.Is(err, models.ErrURLDeleted) {
			c.Status(http.StatusGone)
			return
		}
		c.Status(http.StatusBadRequest)
		return
	}
	c.Header("Location", originURL)
	c.Status(http.StatusTemporaryRedirect)
}

// GetJSONShortURL converts a long URL to its shortened version using JSON input.
// Expects a JSON object with a 'URL' field in the request body.
// Returns a JSON object containing the shortened URL on success.
// Sends HTTP status 400 Bad Request for malformed JSON or invalid URL,
// or HTTP status 409 Conflict if the URL is already shortened.
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

// GetBatchShortURL converts multiple URLs to their shortened versions in batch.
// Expects a JSON array of URL objects in the request body.
// Returns a JSON array of objects containing original and shortened URLs.
// On failure, sends HTTP status 500 Internal Server Error.
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

// GetUserURLS retrieves all URLs associated with the current user.
// Does not require any parameters; user identification is from the context.
// Returns a JSON array of URLs associated with the user.
// If no URLs are found, returns HTTP status 204 No Content.
func (h Handlers) GetUserURLS(c *gin.Context) {
	ctx := c.Request.Context()
	fullShortUserURLS, err := h.service.GetUserURLS(ctx)
	if err != nil {
		c.Status(http.StatusNoContent)
		return
	}
	c.JSON(http.StatusOK, fullShortUserURLS)
}

// DelUserURLS marks specified URLs as deleted asynchronously.
// Expects a JSON array of URL IDs to delete in the request body.
// Acknowledges the deletion request with HTTP status 202 Accepted.
// Returns HTTP status 400 Bad Request for malformed JSON input.
func (h Handlers) DelUserURLS(c *gin.Context) {
	ctx := c.Request.Context()
	var URLSToDel []string
	if err := c.ShouldBindJSON(&URLSToDel); err != nil {
		logrus.Error(err)
		c.Status(http.StatusBadRequest)
		return
	}
	c.Status(http.StatusAccepted)
	h.service.AsyncDeleteUserURLs(ctx, URLSToDel)

}

// PingDB checks the database connection.
// Does not require any parameters.
// Returns HTTP status 200 OK if the database connection is alive.
// On database connection failure, returns HTTP status 500 Internal Server Error.
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
	logrus.Info("DB connection pool is empty")
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

// MiddlewareLogging provides a logging middleware for Gin.
// It logs details about each request including the URL, method, response status, duration, and size.
// This middleware is useful for monitoring and debugging purposes.
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

// MiddlewareCompress provides a compression middleware using gzip.
// It checks the 'Accept-Encoding' header of incoming requests and applies gzip compression if applicable.
// This middleware optimizes response size and speed, improving overall performance.
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

// MiddlewareAuthPublic provides authentication middleware for public routes.
// It manages user tokens, generating new tokens if necessary, and adds user ID to the context.
// This middleware is useful for routes that require user identification but not strict authentication.
func (h Handlers) MiddlewareAuthPublic() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string
		var err error
		var userID uuid.UUID

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

		ctx := context.WithValue(c.Request.Context(), models.UserIDKey, userID)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// MiddlewareAuthPrivate provides authentication middleware for private routes.
// It checks the user token and only allows access if the token is valid.
// This middleware ensures that only authenticated users can access certain routes.
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
		ctx := context.WithValue(c.Request.Context(), models.UserIDKey, userID)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
