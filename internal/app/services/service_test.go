package services

import (
	"errors"
	"fmt"
	"github.com/DenisKhanov/shorterURL/internal/app/services/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewService(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockRepository(ctrl)
	baseURL := "http://localhost:8080"
	service := NewServices(mockRepo, baseURL)
	if service.repository != mockRepo {
		t.Errorf("Expected repository to be set, got %v", service.repository)
	}
}

func TestServices_GetShortURL(t *testing.T) {
	tests := []struct {
		name             string
		originalURL      string
		expectedShortURL string
		existsRepo       bool
		mockSetup        func(mockRepo *mocks.MockRepository)
	}{
		{
			name:             "ShortURL found in repository",
			originalURL:      "http://original.url",
			expectedShortURL: "http://localhost:8080/A4UUE",
			existsRepo:       true,
			mockSetup: func(mockRepo *mocks.MockRepository) {
				mockRepo.EXPECT().GetShortURLFromDB("http://original.url").Return("A4UUE", true).AnyTimes()
			},
		},
		{
			name:             "ShortURL not found in repository",
			originalURL:      "http://original.url",
			expectedShortURL: "http://localhost:8080/A4UUE",
			existsRepo:       false,
			mockSetup: func(mockRepo *mocks.MockRepository) {
				mockRepo.EXPECT().GetShortURLFromDB("http://original.url").Return("", false).AnyTimes()
				mockRepo.EXPECT().GetIDFromDB().Return(214134122).AnyTimes()
				mockRepo.EXPECT().StoreURLSInDB("http://original.url", "A4UUE").Return(nil).AnyTimes()
			},
		},
		{
			name: "ShortURL not found in repository " +
				"and failed to save",
			originalURL:      "http://original.url",
			expectedShortURL: "",
			existsRepo:       false,
			mockSetup: func(mockRepo *mocks.MockRepository) {
				mockRepo.EXPECT().GetShortURLFromDB("http://original.url").Return("", false).AnyTimes()
				mockRepo.EXPECT().GetIDFromDB().Return(214134122).AnyTimes()
				mockRepo.EXPECT().StoreURLSInDB("http://original.url", "A4UUE").Return(errors.New("error saving shortUrl")).AnyTimes()
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
		existsRepo          bool
		mockSetup           func(mockRepo *mocks.MockRepository)
	}{
		{
			name:                "OriginalURL found in repository",
			shortURL:            "A4UUE",
			expectedOriginalURL: "http://original.url",
			existsRepo:          true,
			mockSetup: func(mockRepo *mocks.MockRepository) {
				mockRepo.EXPECT().GetOriginalURLFromDB("A4UUE").Return("http://original.url", true).AnyTimes()
			},
		},
		{
			name:                "OriginalURL not found in repository",
			shortURL:            "A4UUE",
			expectedOriginalURL: "",
			existsRepo:          true,
			mockSetup: func(mockRepo *mocks.MockRepository) {
				mockRepo.EXPECT().GetOriginalURLFromDB("A4UUE").Return("", false).AnyTimes()
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
				assert.EqualError(t, err, fmt.Sprintf("%s/%s not found", service.baseURL, tt.shortURL))
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedOriginalURL, result)
		})
	}
}

func TestBase62Encode(t *testing.T) {
	testCases := []struct {
		input    int
		expected string
	}{
		{0, "0"},
		{10000, "Ib2"},
		{12312313, "hzep"},
		{214134121, "94UUE"},
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
