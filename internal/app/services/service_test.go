package services

import (
	"errors"
	"github.com/DenisKhanov/shorterURL/internal/app/services/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestNewService(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockRepository(ctrl)
	mockEncoder := mocks.NewMockEncoder(ctrl)
	baseURL := "http://localhost:8080"
	service := NewServices(mockRepo, mockEncoder, baseURL)
	if service.repository != mockRepo {
		t.Errorf("Expected repository to be set, got %v", service.repository)
	}
	if service.encoder != mockEncoder {
		t.Errorf("Expected encoder to be set, got %v", service.encoder)
	}
}

func TestServices_GetShortURL(t *testing.T) {

	tests := []struct {
		name             string
		originalURL      string
		expectedShortURL string
		mockSetup        func(mockRepo *mocks.MockRepository, mockEncoder *mocks.MockEncoder)
	}{
		{
			name:             "ShortURL found in repository",
			originalURL:      "http://original.url",
			expectedShortURL: "http://localhost:8080/shortURL",
			mockSetup: func(mockRepo *mocks.MockRepository, mockEncoder *mocks.MockEncoder) {
				mockRepo.EXPECT().GetShortURLFromDB("http://original.url").Return("shortURL", nil).AnyTimes()
			},
		},
		{
			name:             "ShortURL not found in repository",
			originalURL:      "http://original.url",
			expectedShortURL: "http://localhost:8080/shortURL",
			mockSetup: func(mockRepo *mocks.MockRepository, mockEncoder *mocks.MockEncoder) {
				mockRepo.EXPECT().GetShortURLFromDB("http://original.url").Return("", errors.New("short URL not found")).AnyTimes()
				mockEncoder.EXPECT().CryptoBase62Encode().Return("shortURL").AnyTimes()
				mockRepo.EXPECT().StoreURLSInDB("http://original.url", "shortURL").Return(nil).AnyTimes()
			},
		},
		{
			name: "ShortURL not found in repository " +
				"and failed to save",
			originalURL:      "http://original.url",
			expectedShortURL: "",
			mockSetup: func(mockRepo *mocks.MockRepository, mockEncoder *mocks.MockEncoder) {
				mockRepo.EXPECT().GetShortURLFromDB("http://original.url").Return("", errors.New("short URL not found")).AnyTimes()
				mockEncoder.EXPECT().CryptoBase62Encode().Return("shortURL").AnyTimes()
				mockRepo.EXPECT().StoreURLSInDB("http://original.url", "shortURL").Return(errors.New("error saving shortUrl")).AnyTimes()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRepo := mocks.NewMockRepository(ctrl)
			mockEncoder := mocks.NewMockEncoder(ctrl)
			tt.mockSetup(mockRepo, mockEncoder)
			service := Services{repository: mockRepo, encoder: mockEncoder, baseURL: "http://localhost:8080"}
			result, err := service.GetShortURL(tt.originalURL)
			if tt.name == "ShortURL not found in repository "+
				"and failed to save" {
				assert.EqualError(t, err, "error saving shortUrl")
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedShortURL, result)
		})
	}

}

func TestServices_GetOriginalURL(t *testing.T) {
	tests := []struct {
		name                string
		shortURL            string
		expectedOriginalURL string
		mockSetup           func(mockRepo *mocks.MockRepository)
	}{
		{
			name:                "OriginalURL found in repository",
			shortURL:            "shortURL",
			expectedOriginalURL: "http://original.url",
			mockSetup: func(mockRepo *mocks.MockRepository) {
				mockRepo.EXPECT().GetOriginalURLFromDB("shortURL").Return("http://original.url", nil).AnyTimes()
			},
		},
		{
			name:                "OriginalURL not found in repository",
			shortURL:            "shortURL",
			expectedOriginalURL: "",
			mockSetup: func(mockRepo *mocks.MockRepository) {
				mockRepo.EXPECT().GetOriginalURLFromDB("shortURL").Return("", errors.New("original URL not found")).AnyTimes()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRepo := mocks.NewMockRepository(ctrl)
			tt.mockSetup(mockRepo)
			service := Services{repository: mockRepo, baseURL: "http://localhost:8080"}
			result, err := service.GetOriginalURL(tt.shortURL)
			if tt.name == "OriginalURL not found in repository" {
				assert.EqualError(t, err, "original URL not found")
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedOriginalURL, result)
		})
	}
}

func TestCryptoBase62Encode(t *testing.T) {
	service := Services{}

	encoded := service.CryptoBase62Encode()

	if len(encoded) > 7 {
		t.Errorf("expected length <= 7, got %d", len(encoded))
	}

	const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	for _, char := range encoded {
		if !strings.ContainsRune(base62Chars, char) {
			t.Errorf("invalid character %c in encoded string", char)
		}
	}
}
