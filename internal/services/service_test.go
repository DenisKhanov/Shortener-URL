package services

import (
	"context"
	"errors"
	"github.com/DenisKhanov/shorterURL/internal/models"
	"github.com/DenisKhanov/shorterURL/internal/services/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestNewService(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockRepository(ctrl)
	mockEncoder := mocks.NewMockEncoder(ctrl)
	baseURL := "http://localhost:8080"
	service := NewShortURLServices(mockRepo, mockEncoder, baseURL)
	if service.repository != mockRepo {
		t.Errorf("Expected repository to be set, got %v", service.repository)
	}
	if service.encoder != mockEncoder {
		t.Errorf("Expected encoder to be set, got %v", service.encoder)
	}
}
func TestGetBatchShortURL(t *testing.T) {
	tests := []struct {
		name              string
		batchURLRequests  []models.URLRequest
		expectedResponses []models.URLResponse
		repositoryError   error
		storeError        error
		mockSetup         func(mockRepo *mocks.MockRepository, mockEncoder *mocks.MockEncoder)
	}{
		{
			name: "All URLs exist in the database",
			batchURLRequests: []models.URLRequest{
				{CorrelationID: "1", OriginalURL: "http://example1.com"},
				{CorrelationID: "2", OriginalURL: "http://example2.com"},
			},
			expectedResponses: []models.URLResponse{
				{CorrelationID: "1", ShortURL: "http://localhost:8080/short1"},
				{CorrelationID: "2", ShortURL: "http://localhost:8080/short2"},
			},
			repositoryError: nil,
			storeError:      nil,
			mockSetup: func(mockRepo *mocks.MockRepository, mockEncoder *mocks.MockEncoder) {
				shortsURL := map[string]string{
					"http://example1.com": "short1",
					"http://example2.com": "short2",
				}
				mockRepo.EXPECT().GetShortBatchURLFromDB(gomock.Any(), gomock.Eq([]models.URLRequest{
					{CorrelationID: "1", OriginalURL: "http://example1.com"},
					{CorrelationID: "2", OriginalURL: "http://example2.com"},
				})).Return(shortsURL, nil).AnyTimes()
				mockRepo.EXPECT().StoreBatchURLInDB(gomock.Any(), gomock.Any()).Return(nil).Times(1)
			},
		},
		{
			name: "Some URLs need to be generated and stored",
			batchURLRequests: []models.URLRequest{
				{CorrelationID: "3", OriginalURL: "http://example3.com"},
				{CorrelationID: "4", OriginalURL: "http://example4.com"},
			},
			expectedResponses: []models.URLResponse{
				{CorrelationID: "3", ShortURL: "http://localhost:8080/short3"},
				{CorrelationID: "4", ShortURL: "http://localhost:8080/short4"},
			},
			repositoryError: nil,
			storeError:      nil,
			mockSetup: func(mockRepo *mocks.MockRepository, mockEncoder *mocks.MockEncoder) {
				shortsURL := map[string]string{}
				mockRepo.EXPECT().GetShortBatchURLFromDB(gomock.Any(), gomock.Eq([]models.URLRequest{
					{CorrelationID: "3", OriginalURL: "http://example3.com"},
					{CorrelationID: "4", OriginalURL: "http://example4.com"},
				})).Return(shortsURL, nil).AnyTimes()
				mockEncoder.EXPECT().CryptoBase62Encode().Return("short3").AnyTimes().Times(1)
				mockEncoder.EXPECT().CryptoBase62Encode().Return("short4").AnyTimes().Times(1)
				mockRepo.EXPECT().StoreBatchURLInDB(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			},
		},
		{
			name: "Error retrieving URLs from the database",
			batchURLRequests: []models.URLRequest{
				{CorrelationID: "5", OriginalURL: "http://example5.com"},
			},
			expectedResponses: nil,
			repositoryError:   errors.New("database error"),
			storeError:        nil,
			mockSetup: func(mockRepo *mocks.MockRepository, mockEncoder *mocks.MockEncoder) {
				mockRepo.EXPECT().GetShortBatchURLFromDB(gomock.Any(), gomock.Any()).Return(nil, errors.New("database error")).AnyTimes()
			},
		},
		{
			name: "Error storing URLs in the database",
			batchURLRequests: []models.URLRequest{
				{CorrelationID: "6", OriginalURL: "http://example6.com"},
			},
			expectedResponses: nil,
			repositoryError:   nil,
			storeError:        errors.New("storage error"),
			mockSetup: func(mockRepo *mocks.MockRepository, mockEncoder *mocks.MockEncoder) {
				shortsURL := map[string]string{}
				mockRepo.EXPECT().GetShortBatchURLFromDB(gomock.Any(), gomock.Any()).Return(shortsURL, nil).AnyTimes()
				mockEncoder.EXPECT().CryptoBase62Encode().Return("short6").AnyTimes()
				mockRepo.EXPECT().StoreBatchURLInDB(gomock.Any(), gomock.Eq(map[string]string{
					"short6": "http://example6.com",
				})).Return(errors.New("storage error")).AnyTimes()
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

			service := ShortURLServices{repository: mockRepo, encoder: mockEncoder, baseURL: "http://localhost:8080"}
			responses, err := service.GetBatchShortURL(context.Background(), tt.batchURLRequests)

			if tt.expectedResponses != nil {
				assert.NotNil(t, responses)
				assert.Equal(t, tt.expectedResponses, responses)
			} else {
				assert.Nil(t, responses)
			}

			if tt.repositoryError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.repositoryError.Error())
			} else if tt.storeError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.storeError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
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
				mockRepo.EXPECT().GetShortURLFromDB(gomock.Any(), "http://original.url").Return("shortURL", nil).AnyTimes()
			},
		},
		{
			name:             "ShortURL not found in repository",
			originalURL:      "http://original.url",
			expectedShortURL: "http://localhost:8080/shortURL",
			mockSetup: func(mockRepo *mocks.MockRepository, mockEncoder *mocks.MockEncoder) {
				mockRepo.EXPECT().GetShortURLFromDB(context.Background(), "http://original.url").Return("", errors.New("short URL not found")).AnyTimes()
				mockEncoder.EXPECT().CryptoBase62Encode().Return("shortURL").AnyTimes()
				mockRepo.EXPECT().StoreURLInDB(gomock.Any(), "http://original.url", "shortURL").Return(nil).AnyTimes()
			},
		},
		{
			name: "ShortURL not found in repository " +
				"and failed to save",
			originalURL:      "http://original.url",
			expectedShortURL: "",
			mockSetup: func(mockRepo *mocks.MockRepository, mockEncoder *mocks.MockEncoder) {
				mockRepo.EXPECT().GetShortURLFromDB(context.Background(), "http://original.url").Return("", errors.New("short URL not found")).AnyTimes()
				mockEncoder.EXPECT().CryptoBase62Encode().Return("shortURL").AnyTimes()
				mockRepo.EXPECT().StoreURLInDB(gomock.Any(), "http://original.url", "shortURL").Return(errors.New("error saving shortUrl")).AnyTimes()
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
			service := ShortURLServices{repository: mockRepo, encoder: mockEncoder, baseURL: "http://localhost:8080"}
			result, err := service.GetShortURL(context.Background(), tt.originalURL)
			if tt.name == "ShortURL found in repository" {
				assert.Equal(t, tt.expectedShortURL, result)
				assert.EqualError(t, err, "short URL found in database")
			} else {
				if tt.name == "ShortURL not found in repository "+
					"and failed to save" {
					assert.EqualError(t, err, "error saving shortUrl")
				} else {
					assert.NoError(t, err)
				}
			}

			assert.Equal(t, tt.expectedShortURL, result)
		})
	}

}

func TestGetUserURLS(t *testing.T) {
	// Создаем моки
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockEncoder := mocks.NewMockEncoder(ctrl)

	// Создаем экземпляр сервиса
	shortURLService := NewShortURLServices(mockRepo, mockEncoder, "http://localhost:8080")

	// Подготовим тестовые случаи
	testCases := []struct {
		name           string
		userURLSFromDB []models.URL
		expectedOutput []models.URL
		expectedError  error
	}{
		{
			name: "Success case",
			userURLSFromDB: []models.URL{
				{ShortURL: "short1", OriginalURL: "http://example1.com"},
				{ShortURL: "short2", OriginalURL: "http://example2.com"},
			},
			expectedOutput: []models.URL{
				{ShortURL: "http://localhost:8080/short1", OriginalURL: "http://example1.com"},
				{ShortURL: "http://localhost:8080/short2", OriginalURL: "http://example2.com"},
			},
			expectedError: nil,
		},
	}

	// Пройдемся по каждому тестовому случаю
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Устанавливаем ожидания моков
			mockRepo.EXPECT().GetUserURLSFromDB(gomock.Any()).Return(tc.userURLSFromDB, tc.expectedError).AnyTimes()

			// Вызываем тестируемый метод
			actualOutput, actualError := shortURLService.GetUserURLS(context.Background())

			// Проверяем результаты
			assert.Equal(t, tc.expectedError, actualError)

			if tc.expectedError == nil {
				assert.Equal(t, tc.expectedOutput, actualOutput)
			}
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
				mockRepo.EXPECT().GetOriginalURLFromDB(gomock.Any(), "shortURL").Return("http://original.url", nil).AnyTimes()
			},
		},
		{
			name:                "OriginalURL not found in repository",
			shortURL:            "shortURL",
			expectedOriginalURL: "",
			mockSetup: func(mockRepo *mocks.MockRepository) {
				mockRepo.EXPECT().GetOriginalURLFromDB(gomock.Any(), "shortURL").Return("", errors.New("original URL not found")).AnyTimes()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRepo := mocks.NewMockRepository(ctrl)
			tt.mockSetup(mockRepo)
			service := ShortURLServices{repository: mockRepo, baseURL: "http://localhost:8080"}
			result, err := service.GetOriginalURL(context.Background(), tt.shortURL)
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
	service := ShortURLServices{}

	encoded := service.CryptoBase62Encode()

	if len(encoded) > 8 {
		t.Errorf("expected length <= 8, got %d", len(encoded))
	}

	const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	for _, char := range encoded {
		if !strings.ContainsRune(base62Chars, char) {
			t.Errorf("invalid character %c in encoded string", char)
		}
	}
}

func TestAsyncDeleteUserURLs(t *testing.T) {
	// Создаем моки
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockEncoder := mocks.NewMockEncoder(ctrl)

	// Создаем экземпляр сервиса
	shortURLService := NewShortURLServices(mockRepo, mockEncoder, "http://localhost:8080")

	// Подготовим тестовые случаи
	testCases := []struct {
		name          string
		URLsToDelete  []string
		expectedError error
	}{
		{
			name:          "Success case",
			URLsToDelete:  []string{"short1", "short2"},
			expectedError: nil,
		},
	}

	// Пройдемся по каждому тестовому случаю
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Устанавливаем ожидания моков
			mockRepo.EXPECT().MarkURLsAsDeleted(gomock.Any(), gomock.Eq(tc.URLsToDelete)).Return(tc.expectedError).AnyTimes()

			// Вызываем тестируемый метод
			shortURLService.AsyncDeleteUserURLs(context.Background(), tc.URLsToDelete)

			// Ждем некоторое время для завершения асинхронной операции (здесь можно использовать библиотеку для ожидания)
			time.Sleep(time.Second) // Пример использования time.Sleep, лучше использовать специализированные библиотеки для ожидания

			// Проверяем, что ошибок не было
			// Мы не можем проверить успешное выполнение асинхронной операции напрямую, поэтому предполагаем, что если тест завершился без ошибок, то все в порядке.
		})
	}
}

func BenchmarkShortURLServices_CryptoBase62Encode(b *testing.B) {
	service := ShortURLServices{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.CryptoBase62Encode()
	}
}
