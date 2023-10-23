package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/DenisKhanov/shorterURL/internal/app/handlers/mocks"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewHandlers(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockService := mocks.NewMockService(ctrl)
	handlers := NewHandlers(mockService)

	if handlers.service != mockService {
		t.Errorf("Expected service to be set, got %v", handlers.service)
	}
}

func TestHandlers_PostURL(t *testing.T) {

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
				mockService.EXPECT().GetShortURL("http://original.url").Return("94UUE", nil).AnyTimes()
			},
		},
		{
			name:             "POST not valid URL",
			inputURL:         "original.url",
			expectedShortURL: "",
			expectedStatus:   http.StatusBadRequest,
			mockSetup: func(mockService *mocks.MockService) {
				mockService.EXPECT().GetShortURL("original.url").Return("", nil).AnyTimes()
			},
		}, {
			name:             "POST service get error",
			inputURL:         "http://original.url",
			expectedShortURL: "",
			expectedStatus:   http.StatusBadRequest,
			mockSetup: func(mockService *mocks.MockService) {
				mockService.EXPECT().GetShortURL("http://original.url").Return("", errors.New("some error")).AnyTimes()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockService := mocks.NewMockService(ctrl)
			tt.mockSetup(mockService)
			r := httptest.NewRequest("POST", "/", bytes.NewBufferString(tt.inputURL))
			w := httptest.NewRecorder()
			handler := Handlers{service: mockService}
			http.HandlerFunc(handler.PostURL).ServeHTTP(w, r)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, tt.expectedShortURL, w.Body.String())
		})
	}
}

func TestHandlers_GetURL(t *testing.T) {
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
				mockService.EXPECT().GetOriginalURL("94UUE").Return("http://original.url", nil).AnyTimes()
			},
		},
		{
			name:           "GET service get error",
			shortURL:       "/94UUE",
			expectedStatus: http.StatusBadRequest,
			expectedURL:    "",
			mockSetup: func(mockService *mocks.MockService) {
				mockService.EXPECT().GetOriginalURL("94UUE").Return("", errors.New("some error")).AnyTimes()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockService := mocks.NewMockService(ctrl)
			tt.mockSetup(mockService)
			handler := Handlers{service: mockService}
			r := httptest.NewRequest("GET", tt.shortURL, nil)
			w := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/{id}", handler.GetURL).Methods("GET")
			router.ServeHTTP(w, r)
			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, tt.expectedURL, w.Header().Get("Location"))
		})
	}
}

func TestHandlers_JsonURL(t *testing.T) {

	tests := []struct {
		name           string
		inputJSON      string
		expectedJSON   string
		expectedStatus int
		mockSetup      func(mockService *mocks.MockService)
	}{
		{
			name:           "POST Valid URL",
			inputJSON:      `{"url": "http://original.url"}`,
			expectedJSON:   `{"result": "http://localhost:8080/94UUE"}`,
			expectedStatus: http.StatusCreated,
			mockSetup: func(mockService *mocks.MockService) {
				mockService.EXPECT().GetShortURL("http://original.url").Return("http://localhost:8080/94UUE", nil).AnyTimes()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockService := mocks.NewMockService(ctrl)
			tt.mockSetup(mockService)
			r := httptest.NewRequest("POST", "/api/shorten", bytes.NewBufferString(tt.inputJSON))
			w := httptest.NewRecorder()
			handler := Handlers{service: mockService}
			http.HandlerFunc(handler.JSONURL).ServeHTTP(w, r)

			assert.Equal(t, tt.expectedStatus, w.Code)
			var actual, expected map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &actual)
			json.Unmarshal([]byte(tt.expectedJSON), &expected)
			assert.Equal(t, expected, actual)
		})
	}
}
