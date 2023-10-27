package handlers

import (
	"compress/gzip"
	"encoding/json"
	"github.com/gorilla/mux"
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
}

type URLProcessingResult struct {
	URL    string `json:"url"`
	Result string `json:"result"`
}

func (u URLProcessingResult) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Result string `json:"result"`
	}{
		Result: u.Result,
	})
}

var typeArray = [2]string{"application/json", "text/html"}

type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
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

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
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

func (c *compressReader) Read(data []byte) (n int, err error) {
	return c.zr.Read(data)
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

type Handlers struct {
	service Service
}

func NewHandlers(service Service) *Handlers {
	return &Handlers{
		service: service,
	}
}

func (h Handlers) PostURL(w http.ResponseWriter, r *http.Request) {
	linc, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
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
func (h Handlers) JSONURL(w http.ResponseWriter, r *http.Request) {
	var dataURL URLProcessingResult
	if err := json.NewDecoder(r.Body).Decode(&dataURL); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	result, err := h.service.GetShortURL(dataURL.URL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	dataURL.Result = result

	jsonResult, err := json.Marshal(dataURL)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(jsonResult)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

// logging
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
func (h Handlers) MiddlewareCompress(ha http.Handler) http.Handler {
	compressFn := func(w http.ResponseWriter, r *http.Request) {
		ow := w
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			ha.ServeHTTP(w, r)
			return
		}
		cw := newCompressWriter(w)
		defer cw.Close()
		for _, v := range typeArray {
			if strings.Contains(r.Header.Get("Content-Type"), v) {
				ow = cw
				break
			}
		}
		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = cr
			defer cr.Close()
		}

		ha.ServeHTTP(ow, r)
	}
	return http.HandlerFunc(compressFn)
}
