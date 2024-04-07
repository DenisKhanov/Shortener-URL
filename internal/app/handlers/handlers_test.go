package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/DenisKhanov/shorterURL/internal/app/handlers/mocks"
	"github.com/DenisKhanov/shorterURL/internal/app/models"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlers_GetShortURL(t *testing.T) {

	tests := []struct {
		name             string
		inputURL         string
		expectedShortURL string
		expectedStatus   int
		mockSetup        func(mockService *mocks.MockService)
	}{
		{
			name:             "POST Valid URL",
			inputURL:         "http://original.url",
			expectedShortURL: "94UUE",
			expectedStatus:   http.StatusCreated,
			mockSetup: func(mockService *mocks.MockService) {
				mockService.EXPECT().GetShortURL(gomock.Any(), "http://original.url").Return("94UUE", nil).AnyTimes()
			},
		},
		{
			name:             "POST not valid URL",
			inputURL:         "original.url",
			expectedShortURL: "",
			expectedStatus:   http.StatusBadRequest,
			mockSetup: func(mockService *mocks.MockService) {
				mockService.EXPECT().GetShortURL(gomock.Any(), "original.url").Return("", nil).AnyTimes()
			},
		}, {
			name:             "POST service get error",
			inputURL:         "http://original.url",
			expectedShortURL: "{\"error\":\"some error\"}",
			expectedStatus:   http.StatusBadRequest,
			mockSetup: func(mockService *mocks.MockService) {
				mockService.EXPECT().GetShortURL(gomock.Any(), "http://original.url").Return("", errors.New("some error")).AnyTimes()
			},
		},
		{
			name:             "POST service get error models.ErrURLFound",
			inputURL:         "http://original.url",
			expectedShortURL: "",
			expectedStatus:   http.StatusConflict,
			mockSetup: func(mockService *mocks.MockService) {
				mockService.EXPECT().GetShortURL(gomock.Any(), "http://original.url").Return("", models.ErrURLFound).AnyTimes()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создание Gin контекста
			gin.SetMode(gin.TestMode)
			r := gin.Default()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockService := mocks.NewMockService(ctrl)
			tt.mockSetup(mockService)
			handler := Handlers{service: mockService}
			r.POST("/", handler.GetShortURL)

			// Создание HTTP запроса и рекордера
			req := httptest.NewRequest("POST", "/", bytes.NewBufferString(tt.inputURL))
			w := httptest.NewRecorder()

			// Выполнение запроса через Gin
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, tt.expectedShortURL, w.Body.String())
		})
	}
}

func TestHandlers_GetOriginalURL(t *testing.T) {
	tests := []struct {
		name           string
		shortURL       string
		expectedStatus int
		expectedURL    string
		mockSetup      func(mockService *mocks.MockService)
	}{
		{
			name:           "GET valid shortURL",
			shortURL:       "/94UUE",
			expectedStatus: http.StatusTemporaryRedirect,
			expectedURL:    "http://original.url",
			mockSetup: func(mockService *mocks.MockService) {
				mockService.EXPECT().GetOriginalURL(gomock.Any(), "94UUE").Return("http://original.url", nil).AnyTimes()
			},
		},
		{
			name:           "GET service get error",
			shortURL:       "/94UUE",
			expectedStatus: http.StatusBadRequest,
			expectedURL:    "",
			mockSetup: func(mockService *mocks.MockService) {
				mockService.EXPECT().GetOriginalURL(gomock.Any(), "94UUE").Return("", errors.New("some error")).AnyTimes()
			},
		},
		{
			name:           "GET service get error models.ErrURLDeleted",
			shortURL:       "/94UUE",
			expectedStatus: http.StatusGone,
			expectedURL:    "",
			mockSetup: func(mockService *mocks.MockService) {
				mockService.EXPECT().GetOriginalURL(gomock.Any(), "94UUE").Return("", models.ErrURLDeleted).AnyTimes()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockService := mocks.NewMockService(ctrl)
			tt.mockSetup(mockService)

			// Создание Gin роутера
			gin.SetMode(gin.TestMode)
			router := gin.Default()
			handler := Handlers{service: mockService}
			router.GET("/:id", handler.GetOriginalURL)

			// Создание HTTP запроса и рекордера
			req := httptest.NewRequest("GET", tt.shortURL, nil)
			w := httptest.NewRecorder()

			// Выполнение запроса через Gin
			router.ServeHTTP(w, req)

			// Проверки
			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, tt.expectedURL, w.Header().Get("Location"))
		})
	}
}

func TestHandlers_GetJSONShortURL(t *testing.T) {

	tests := []struct {
		name           string
		inputJSON      string
		expectedJSON   string
		expectedStatus int
		expectedError  error
		mockSetup      func(mockService *mocks.MockService)
	}{
		{
			name:           "POST Valid URL and shortURL not found in database",
			inputJSON:      `{"url": "http://original.url"}`,
			expectedJSON:   `{"result": "http://localhost:8080/94UUE"}`,
			expectedStatus: http.StatusCreated,
			mockSetup: func(mockService *mocks.MockService) {
				mockService.EXPECT().GetShortURL(gomock.Any(), "http://original.url").Return("http://localhost:8080/94UUE", nil).AnyTimes()
			},
		},
		{
			name:           "POST Valid URL and shortURL found in database",
			inputJSON:      `{"url": "http://original.url"}`,
			expectedJSON:   `{"result": "http://localhost:8080/94UUE"}`,
			expectedStatus: http.StatusConflict,
			expectedError:  models.ErrURLFound,
			mockSetup: func(mockService *mocks.MockService) {
				mockService.EXPECT().GetShortURL(gomock.Any(), "http://original.url").Return("http://localhost:8080/94UUE", models.ErrURLFound).AnyTimes()
			},
		},
		{
			name:           "POST Valid URL and shortURL found in database",
			inputJSON:      `{"url": "http://original.url"}`,
			expectedJSON:   `{"result": "http://localhost:8080/94UUE"}`,
			expectedStatus: http.StatusConflict,
			expectedError:  models.ErrURLFound,
			mockSetup: func(mockService *mocks.MockService) {
				mockService.EXPECT().GetShortURL(gomock.Any(), "http://original.url").Return("http://localhost:8080/94UUE", models.ErrURLFound).AnyTimes()
			},
		},
		{
			name:           "POST Invalid URL",
			inputJSON:      `{"url": "invalid-url"}`,
			expectedJSON:   `{"error":"any error"}`,
			expectedStatus: http.StatusBadRequest,
			mockSetup: func(mockService *mocks.MockService) {
				mockService.EXPECT().GetShortURL(gomock.Any(), "invalid-url").Return("", errors.New("any error")).AnyTimes()
			},
		},
		{
			name:           "POST Invalid JSON",
			inputJSON:      `invalid_json`,
			expectedJSON:   `{"error": "JSON"}`,
			expectedStatus: http.StatusBadRequest,
			mockSetup: func(mockService *mocks.MockService) {
				// В этом случае mockService не должен вызывать GetBatchShortURL
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создание Gin контекста
			gin.SetMode(gin.TestMode)
			r := gin.Default()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockService := mocks.NewMockService(ctrl)
			tt.mockSetup(mockService)
			handler := Handlers{service: mockService}
			r.POST("/api/shorten", handler.GetJSONShortURL)

			// Создание HTTP запроса и рекордера
			req := httptest.NewRequest("POST", "/api/shorten", bytes.NewBufferString(tt.inputJSON))
			w := httptest.NewRecorder()

			// Выполнение запроса через Gin
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			var actual, expected map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &actual)
			json.Unmarshal([]byte(tt.expectedJSON), &expected)
			assert.Equal(t, expected, actual)
		})
	}
}

func TestHandlers_GetBatchShortURL(t *testing.T) {
	tests := []struct {
		name           string
		inputJSON      string
		expectedJSON   string
		expectedStatus int
		mockSetup      func(mockService *mocks.MockService)
	}{
		{
			name:           "POST Valid Batch URL",
			inputJSON:      `[{"correlation_id": "id1", "original_url": "http://original1.url"}, {"correlation_id": "id2", "original_url": "http://original2.url"}]`,
			expectedJSON:   `[{"correlation_id": "id1", "short_url": "http://localhost:8080/short1"}, {"correlation_id": "id2", "short_url": "http://localhost:8080/short2"}]`,
			expectedStatus: http.StatusCreated,
			mockSetup: func(mockService *mocks.MockService) {
				mockService.EXPECT().GetBatchShortURL(gomock.Any(), gomock.Any()).Return([]models.URLResponse{
					{CorrelationID: "id1", ShortURL: "http://localhost:8080/short1"},
					{CorrelationID: "id2", ShortURL: "http://localhost:8080/short2"},
				}, nil).AnyTimes()
			},
		},
		{
			name:           "POST Invalid Batch JSON",
			inputJSON:      `invalid_json`,
			expectedStatus: http.StatusBadRequest,
			mockSetup: func(mockService *mocks.MockService) {
				// В этом случае mockService не должен вызывать GetBatchShortURL
			},
		},
		{
			name:           "POST Error from Service",
			inputJSON:      `[{"correlation_id": "id1", "original_url": "http://original1.url"}, {"correlation_id": "id2", "original_url": "http://original2.url"}]`,
			expectedStatus: http.StatusInternalServerError,
			mockSetup: func(mockService *mocks.MockService) {
				mockService.EXPECT().GetBatchShortURL(gomock.Any(), gomock.Any()).Return(nil, errors.New("service error")).AnyTimes()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создание Gin контекста
			gin.SetMode(gin.TestMode)
			r := gin.Default()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockService := mocks.NewMockService(ctrl)
			tt.mockSetup(mockService)
			handler := Handlers{service: mockService}
			r.POST("/api/batch-shorten", handler.GetBatchShortURL)

			// Создание HTTP запроса и рекордера
			req := httptest.NewRequest("POST", "/api/batch-shorten", bytes.NewBufferString(tt.inputJSON))
			w := httptest.NewRecorder()

			// Выполнение запроса через Gin
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedJSON != "" {
				var actual, expected []models.URLResponse
				json.Unmarshal(w.Body.Bytes(), &actual)
				json.Unmarshal([]byte(tt.expectedJSON), &expected)
				assert.Equal(t, expected, actual)
			}
		})
	}
}

func TestHandlers_GetUserURLS(t *testing.T) {
	tests := []struct {
		name           string
		expectedJSON   string
		expectedStatus int
		mockSetup      func(mockService *mocks.MockService)
	}{
		{
			name:           "User has URLs",
			expectedJSON:   `[{"original_url": "http://original1.url", "short_url": "http://localhost:8080/short1"}, {"original_url": "http://original2.url", "short_url": "http://localhost:8080/short2"}]`,
			expectedStatus: http.StatusOK,
			mockSetup: func(mockService *mocks.MockService) {
				mockService.EXPECT().GetUserURLS(gomock.Any()).Return([]models.URL{
					{OriginalURL: "http://original1.url", ShortURL: "http://localhost:8080/short1"},
					{OriginalURL: "http://original2.url", ShortURL: "http://localhost:8080/short2"},
				}, nil).AnyTimes()
			},
		},
		{
			name:           "User has no URLs",
			expectedJSON:   "",
			expectedStatus: http.StatusNoContent,
			mockSetup: func(mockService *mocks.MockService) {
				mockService.EXPECT().GetUserURLS(gomock.Any()).Return(nil, errors.New("any error")).AnyTimes()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создание Gin контекста
			gin.SetMode(gin.TestMode)
			r := gin.Default()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockService := mocks.NewMockService(ctrl)
			tt.mockSetup(mockService)
			handler := Handlers{service: mockService}
			r.GET("/api/user-urls", handler.GetUserURLS)

			// Создание HTTP запроса и рекордера
			req := httptest.NewRequest("GET", "/api/user-urls", nil)
			w := httptest.NewRecorder()

			// Выполнение запроса через Gin
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedJSON != "" {
				var actual, expected []models.URL
				json.Unmarshal(w.Body.Bytes(), &actual)
				json.Unmarshal([]byte(tt.expectedJSON), &expected)
				assert.Equal(t, expected, actual)
			}
		})
	}
}

func TestHandlers_DelUserURLS(t *testing.T) {
	tests := []struct {
		name           string
		inputJSON      string
		expectedStatus int
		mockSetup      func(mockService *mocks.MockService)
	}{
		{
			name:           "Valid JSON Request",
			inputJSON:      `["id1", "id2"]`,
			expectedStatus: http.StatusAccepted,
			mockSetup: func(mockService *mocks.MockService) {
				mockService.EXPECT().AsyncDeleteUserURLs(gomock.Any(), []string{"id1", "id2"}).AnyTimes()
			},
		},
		{
			name:           "Invalid JSON Request",
			inputJSON:      `invalid_json`,
			expectedStatus: http.StatusBadRequest,
			mockSetup: func(mockService *mocks.MockService) {
				// В этом случае mockService не должен вызывать AsyncDeleteUserURLs
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создание Gin контекста
			gin.SetMode(gin.TestMode)
			r := gin.Default()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockService := mocks.NewMockService(ctrl)
			tt.mockSetup(mockService)
			handler := Handlers{service: mockService}
			r.DELETE("/api/delete-urls", handler.DelUserURLS)

			// Создание HTTP запроса и рекордера
			req := httptest.NewRequest("DELETE", "/api/delete-urls", bytes.NewBufferString(tt.inputJSON))
			w := httptest.NewRecorder()

			// Выполнение запроса через Gin
			r.ServeHTTP(w, req)

			// Проверки остаются теми же
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestHandlers_MiddlewareLogging(t *testing.T) {
	// Создание контроллера gomock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создание мок-объекта
	mockService := mocks.NewMockService(ctrl)

	// Создание Handlers с использованием мок-объектов
	handler := Handlers{
		service: mockService,
	}

	// Создание Gin контекста
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// Создание буфера для логов
	var logsBuffer bytes.Buffer
	logrus.SetOutput(&logsBuffer)

	// Регистрация middleware
	r.Use(handler.MiddlewareLogging())

	// Обработчик для теста
	r.GET("/api/test-middleware", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// Симуляция запроса
	req := httptest.NewRequest("GET", "/api/test-middleware", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Проверка логов
	logs := logsBuffer.String()
	assert.Contains(t, logs, "/api/test-middleware")
	assert.Contains(t, logs, "GET")
	assert.Contains(t, logs, "200")
	assert.Contains(t, logs, "Обработан запрос")
}

func TestHandlers_MiddlewareCompress(t *testing.T) {
	tests := []struct {
		name                    string
		acceptEncodingHeader    string
		expectedStatus          int
		expectedContentEncoding string
	}{
		{
			name:                    "No compression",
			acceptEncodingHeader:    "identity",
			expectedStatus:          http.StatusOK,
			expectedContentEncoding: "",
		},
		{
			name:                    "Gzip compression",
			acceptEncodingHeader:    "gzip",
			expectedStatus:          http.StatusOK,
			expectedContentEncoding: "gzip",
		},
	}

	// Создание контроллера gomock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создание мок-объекта
	mockService := mocks.NewMockService(ctrl)

	// Создание Handlers с использованием мок-объектов
	handler := Handlers{
		service: mockService,
	}

	// Создание Gin контекста
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// Регистрация middleware
	r.Use(handler.MiddlewareCompress())

	// Обработчик для теста
	r.POST("/api/test-compress", func(c *gin.Context) {
		// Ваш код обработки запроса
		c.Status(http.StatusOK)
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Симуляция запроса
			req := httptest.NewRequest("POST", "/api/test-compress", nil)
			req.Header.Set("Accept-Encoding", tt.acceptEncodingHeader)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			// Проверки
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedContentEncoding != "" {
				assert.Contains(t, w.Header(), "Content-Encoding")
				assert.Equal(t, tt.expectedContentEncoding, w.Header().Get("Content-Encoding"))
			} else {
				assert.NotContains(t, w.Header(), "Content-Encoding")
			}
		})
	}
}
