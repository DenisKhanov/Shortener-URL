package handlers

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"time"
)

//go:generate mockgen -source=handlers.go -destination=mocks/handlers_mock.go -package=mocks
type Service interface {
	GetShortURL(url string) (string, error)
	GetOriginalURL(shortURL string) (string, error)
}

type Handlers struct {
	service Service
}

func NewHandlers(service Service) *Handlers {
	return &Handlers{
		service: service,
	}
}

// PostURL

func (h Handlers) PostURL(w http.ResponseWriter, r *http.Request) {
	linc, _ := io.ReadAll(r.Body)
	lincString := string(linc)
	defer r.Body.Close()

	parsedLinc, err := url.Parse(lincString)
	if err != nil || parsedLinc.Scheme == "" || parsedLinc.Host == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	shortURL, err := h.service.GetShortURL(lincString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))

}
func (h Handlers) GetURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortURL := vars["id"]
	originURL, err := h.service.GetOriginalURL(shortURL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", originURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

type responseData struct {
	status int
	size   int
}
type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
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

func (h Handlers) MiddlewareLogging(ha http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		url := r.RequestURI
		method := r.Method
		duration := time.Since(start)
		ha.ServeHTTP(&lw, r)
		logrus.WithFields(logrus.Fields{
			"url":      url,
			"method":   method,
			"status":   responseData.status,
			"duration": duration,
			"size":     responseData.size,
		}).Info("Обработан запрос")
	}
	return http.HandlerFunc(logFn)
}
