package main

import (
	"errors"
	"github.com/DenisKhanov/shorterURL/internal/app/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestServices_GetShortURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockRepo.EXPECT().GetShortURL("http://original.url").Return("shortURL", true)

	services := NewServices(mockRepo, nil)

	shortURL := services.GetShortURL("http://original.url")

	assert.Equal(t, "http://localhost:8080/shortURL", shortURL)
}

func TestServices_GetOriginURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockRepo.EXPECT().GetOriginalURL("short_url").Return("test_url", true)

	service := Services{storage: mockRepo}

	result, err := service.GetOriginURL("short_url")
	assert.Nil(t, err)
	assert.Equal(t, "test_url", result)

	mockRepo.EXPECT().GetOriginalURL("invalid_url").Return("", false)
	_, err = service.GetOriginURL("invalid_url")
	assert.Equal(t, errors.New("http://localhost:8080/invalid_url not found"), err)
}

func TestBase62Encode(t *testing.T) {
	testCases := []struct {
		input    int
		expected string
	}{
		{10000, "Ib2"},
		{12312313, "hzep"},
		{123123132121, "h77SOA2"},
		{12312312, "gzep"},
		{123123132121, "h77SOA2"},
		{121123122121, "t9H6D82"},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			output := base62Encode(tc.input)
			assert.Equal(t, tc.expected, output)
		})
	}
}
