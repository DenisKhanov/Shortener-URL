package handlers

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"github.com/DenisKhanov/shorterURL/internal/app/models"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

//go:generate mockgen -source=handlers.go -destination=mocks/handlers_mock.go -package=mocks
type Service interface {
	GetShortURL(url string) (string, error)
	GetOriginalURL(shortURL string) (string, error)
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

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}
type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}
type sizeTrackingResponseWriter struct {
	http.ResponseWriter
	size int
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

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}
func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}
func NewHandlers(service Service, DB *pgxpool.Pool) *Handlers {
	return &Handlers{
		service: service,
		DB:      DB,
	}
}

func (h Handlers) GetShortURL(w http.ResponseWriter, r *http.Request) {
	linc, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	r.Body.Close()
	lincString := string(linc)

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
func (h Handlers) GetOriginalURL(w http.ResponseWriter, r *http.Request) {
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
func (h Handlers) GetJSONShortURL(w http.ResponseWriter, r *http.Request) {
	var dataURL URLProcessing
	var dataURLResylt URLProcessingResult
	if err := json.NewDecoder(r.Body).Decode(&dataURL); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	r.Body.Close()

	result, err := h.service.GetShortURL(dataURL.URL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	dataURLResylt.Result = result
	jsonResult, err := json.Marshal(dataURLResylt)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(jsonResult)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
func (h Handlers) GetBatchJSONShortURL(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hendler GetBatchJSONShortURL run")
	var batchURLRequests []models.URLRequest
	if err := json.NewDecoder(r.Body).Decode(&batchURLRequests); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	r.Body.Close()
	batchURLResponses, err := h.service.GetBatchJSONShortURL(batchURLRequests)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonResult, err := json.Marshal(batchURLResponses)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(jsonResult)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (h Handlers) PingDB(w http.ResponseWriter, r *http.Request) {
	if err := h.DB.Ping(context.Background()); err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func (sw *sizeTrackingResponseWriter) Write(data []byte) (int, error) {
	n, err := sw.ResponseWriter.Write(data)
	sw.size += n
	return n, err
}

func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}
func (c *compressWriter) Write(data []byte) (int, error) {
	return c.zw.Write(data)
}
func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}
func (c *compressWriter) Close() error {
	return c.zw.Close()
}
func (c *compressReader) Read(data []byte) (n int, err error) {
	return c.zr.Read(data)
}
func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
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
func (h Handlers) MiddlewareCompress(ha http.Handler) http.Handler {
	compressFn := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			defer cr.Close()
			r.Body = cr
		}
		sw := &sizeTrackingResponseWriter{ResponseWriter: w}
		cw := newCompressWriter(sw)
		defer cw.Close()

		ha.ServeHTTP(cw, r)

		contentType := w.Header().Get("Content-Type")
		if sw.size > 1400 && strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			for _, v := range typeArray {
				if strings.Contains(contentType, v) {
					w.Header().Set("Content-Encoding", "gzip")
					break
				}
			}
		}
	}
	return http.HandlerFunc(compressFn)
}
