package main

import (
	"bytes"
	"github.com/DenisKhanov/shorterURL/internal/app/mocks"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlers_PostURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	mockService.EXPECT().GetShortURL("http://original.url").Return("shortURL").AnyTimes()

	handlers := NewHandlers(mockService, nil)

	request := httptest.NewRequest("POST", "/", bytes.NewBufferString("http://original.url"))
	w := httptest.NewRecorder()

	handlers.PostURL(w, request)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "shortURL", w.Body.String())
}

func TestHandlers_GetURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	mockService.EXPECT().GetOriginURL("shortURL").Return("http://original.url", nil)

	handlers := NewHandlers(mockService, nil)

	request := httptest.NewRequest("GET", "/shortURL", nil)

	w := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/{id}", handlers.GetURL).Methods("GET")
	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
	assert.Equal(t, "http://original.url", w.Header().Get("Location"))
}
