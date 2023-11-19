package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/DenisKhanov/shorterURL/internal/app/handlers/mocks"
	"github.com/DenisKhanov/shorterURL/internal/app/models"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
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
			expectedShortURL: "",
			expectedStatus:   http.StatusBadRequest,
			mockSetup: func(mockService *mocks.MockService) {
				mockService.EXPECT().GetShortURL(gomock.Any(), "http://original.url").Return("", errors.New("some error")).AnyTimes()
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

			// Проверки остаются теми же
			assert.Equal(t, tt.expectedStatus, w.Code)
			var actual, expected map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &actual)
			json.Unmarshal([]byte(tt.expectedJSON), &expected)
			assert.Equal(t, expected, actual)
		})
	}
}
