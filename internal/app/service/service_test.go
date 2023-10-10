package service

import (
	"errors"
	"github.com/DenisKhanov/shorterURL/internal/app/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestServices_GetShortURL_ExistingURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockRepository(ctrl)
	mockRepo.EXPECT().GetShortURL("http://original.url").Return("shortURL", true)

	service := NewServices(mockRepo, nil)
	shortURL, err := service.GetShortURL("http://original.url")

	assert.NoError(t, err)
	assert.Equal(t, "http://localhost:8080/shortURL", shortURL)
}

func TestServices_GetShortURL_NewURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockRepository(gomock.NewController(t))
	mockRepo.EXPECT().GetShortURL("http://original.url").Return("", false)
	mockRepo.EXPECT().GetID().Return(1).AnyTimes()
	mockRepo.EXPECT().StoreURL(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	service := NewServices(mockRepo, nil)
	shortURL, err := service.GetShortURL("http://original.url")
	id := service.storage.GetID()
	assert.Equal(t, 1, id)
	assert.Equal(t, "http://localhost:8080/1", shortURL)
	assert.Nil(t, err)
}
func TestServices_GetShortURL_StoreURLError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockRepository(ctrl)
	mockRepo.EXPECT().GetShortURL("http://original.url").Return("", false)
	mockRepo.EXPECT().GetID().Return(1)
	mockRepo.EXPECT().StoreURL("http://original.url", gomock.Any()).Return(errors.New("storage error"))

	service := NewServices(mockRepo, nil)
	shortURL, err := service.GetShortURL("http://original.url")

	assert.Error(t, err)
	assert.Equal(t, "storage error", err.Error())
	assert.Equal(t, "", shortURL)
}

func TestServices_GetOriginURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockRepo.EXPECT().GetOriginalURL("shortURL").Return("http://original.url", true)

	service := Services{storage: mockRepo}

	result, err := service.GetOriginURL("shortURL")
	assert.Nil(t, err)
	assert.Equal(t, "http://original.url", result)

	mockRepo.EXPECT().GetOriginalURL("invalid_url").Return("", false)
	_, err = service.GetOriginURL("invalid_url")
	assert.Equal(t, errors.New("http://localhost:8080/invalid_url not found"), err)
}

func TestBase62Encode(t *testing.T) {
	testCases := []struct {
		input    int
		expected string
	}{
		{0, "0"},
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
