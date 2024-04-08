package repositories

import (
	"context"
	"fmt"
	"github.com/DenisKhanov/shorterURL/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var UserID, _ = uuid.Parse("e774844b-5895-4b08-b867-50480263f75b")

func TestNewRepository(t *testing.T) {

	tests := []struct {
		name string
		want *URLInMemoryRepo
	}{
		{
			name: "Valid args",
			want: &URLInMemoryRepo{
				shortToOrigURL: map[string]string{"short1": "original1"},
				origToShortURL: map[string]string{"original1": "short1"},
				usersURLS: map[uuid.UUID][]models.URL{
					UserID: {{ShortURL: "short1", OriginalURL: "original1"}},
				},
				batchBuffer:     make([]URLInFileRepo, 0),
				batchCounter:    0,
				batchSize:       100,
				storageFilePath: createTempFilePath(t),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewURLInMemoryRepo(tt.want.storageFilePath), "NewURLInMemoryRepo()")
			defer os.Remove(tt.want.storageFilePath)
		})
	}
}

func TestRepositoryURL_GetOriginalURL(t *testing.T) {
	type fields struct {
		shortToOrigURL map[string]string
		origToShortURL map[string]string
		usersURLS      map[uint32][]models.URL
	}
	type args struct {
		ctx      context.Context
		shortURL string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "valid get original URL",
			fields: fields{
				map[string]string{"short1": "original1"},
				map[string]string{"original1": "short1"},
				map[uint32][]models.URL{123456: {{ShortURL: "short1", OriginalURL: "original1"}}},
			},
			args:    args{ctx: context.Background(), shortURL: "short1"},
			want:    "original1",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &URLInMemoryRepo{
				shortToOrigURL: tt.fields.shortToOrigURL,
				origToShortURL: tt.fields.origToShortURL,
			}
			got, err := d.GetOriginalURL(tt.args.ctx, tt.args.shortURL)
			if !tt.wantErr(t, err, fmt.Sprintf("GetOriginalURL(%v)", tt.args.shortURL)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetOriginalURL(%v)", tt.args.shortURL)
		})
	}
}

func TestRepositoryURL_GetShortURL(t *testing.T) {
	type fields struct {
		shortToOrigURL map[string]string
		origToShortURL map[string]string
		usersURLS      map[uuid.UUID][]models.URL
	}
	type args struct {
		ctx         context.Context
		originalURL string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "valid get short URL",
			fields: fields{
				map[string]string{"short1": "original1"},
				map[string]string{"original1": "short1"},
				map[uuid.UUID][]models.URL{UserID: {{ShortURL: "short1", OriginalURL: "original1"}}},
			},
			args:    args{ctx: context.Background(), originalURL: "original1"},
			want:    "short1",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &URLInMemoryRepo{
				shortToOrigURL: tt.fields.shortToOrigURL,
				origToShortURL: tt.fields.origToShortURL,
			}
			got, err := d.GetShortURL(tt.args.ctx, tt.args.originalURL)
			if !tt.wantErr(t, err, fmt.Sprintf("GetShortURL(%v)", tt.args.originalURL)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetShortURL(%v)", tt.args.originalURL)
		})
	}
}

func TestURLInMemoryRepo_ReadFileToMemoryURL(t *testing.T) {
	type fields struct {
		shortToOrigURL  map[string]string
		origToShortURL  map[string]string
		usersURLS       map[uuid.UUID][]models.URL
		batchBuffer     []URLInFileRepo
		batchCounter    uint8
		batchSize       uint8
		storageFilePath string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "valid read to memory",
			fields: fields{
				shortToOrigURL: map[string]string{"short1": "original1"},
				origToShortURL: map[string]string{"original1": "short1"},
				usersURLS: map[uuid.UUID][]models.URL{
					UserID: {{ShortURL: "short1", OriginalURL: "original1"}},
				},
				batchBuffer:     []URLInFileRepo{},
				batchCounter:    0,
				batchSize:       100,
				storageFilePath: createTempFilePath(t),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &URLInMemoryRepo{
				shortToOrigURL:  tt.fields.shortToOrigURL,
				origToShortURL:  tt.fields.origToShortURL,
				usersURLS:       tt.fields.usersURLS,
				batchBuffer:     tt.fields.batchBuffer,
				batchCounter:    tt.fields.batchCounter,
				batchSize:       tt.fields.batchSize,
				storageFilePath: tt.fields.storageFilePath,
			}
			tt.wantErr(t, m.readFileToMemoryURL(), "ReadFileToMemoryURL()")
			defer os.Remove(tt.fields.storageFilePath)
		})
	}
}

func TestURLInMemoryRepo_SaveBatchToFile(t *testing.T) {
	type fields struct {
		shortToOrigURL  map[string]string
		origToShortURL  map[string]string
		usersURLS       map[uuid.UUID][]models.URL
		batchBuffer     []URLInFileRepo
		batchCounter    uint8
		batchSize       uint8
		storageFilePath string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "valid save batch to file",
			fields: fields{
				shortToOrigURL: map[string]string{"short1": "original1"},
				origToShortURL: map[string]string{"original1": "short1"},
				usersURLS: map[uuid.UUID][]models.URL{
					UserID: {{ShortURL: "short1", OriginalURL: "original1"}},
				},
				batchBuffer: []URLInFileRepo{
					{
						UserID:      UserID,
						ShortURL:    "short1",
						OriginalURL: "original1",
					},
				},
				batchCounter:    1,
				batchSize:       100,
				storageFilePath: createTempFilePath(t),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &URLInMemoryRepo{
				shortToOrigURL:  tt.fields.shortToOrigURL,
				origToShortURL:  tt.fields.origToShortURL,
				usersURLS:       tt.fields.usersURLS,
				batchBuffer:     tt.fields.batchBuffer,
				batchCounter:    tt.fields.batchCounter,
				batchSize:       tt.fields.batchSize,
				storageFilePath: tt.fields.storageFilePath,
			}
			tt.wantErr(t, m.SaveBatchToFile(), "SaveBatchToFile()")
			defer os.Remove(tt.fields.storageFilePath)
		})
	}
}

func TestURLInMemoryRepo_StoreURLS(t *testing.T) {
	type fields struct {
		shortToOrigURL  map[string]string
		origToShortURL  map[string]string
		usersURLS       map[uuid.UUID][]models.URL
		batchBuffer     []URLInFileRepo
		batchCounter    uint8
		batchSize       uint8
		storageFilePath string
	}
	type args struct {
		ctx         context.Context
		originalURL string
		shortURL    string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "valid save URL to memory",
			fields: fields{
				shortToOrigURL: map[string]string{"short2": "original2"},
				origToShortURL: map[string]string{"original2": "short2"},
				usersURLS: map[uuid.UUID][]models.URL{
					UserID: {{ShortURL: "short2", OriginalURL: "original2"}},
				},
				batchBuffer: []URLInFileRepo{
					{
						UserID:      UserID,
						ShortURL:    "short2",
						OriginalURL: "original2",
					},
				},
				batchCounter:    1,
				batchSize:       100,
				storageFilePath: createTempFilePath(t),
			},
			args:    args{ctx: context.Background(), originalURL: "original2", shortURL: "short2"},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &URLInMemoryRepo{
				shortToOrigURL:  tt.fields.shortToOrigURL,
				origToShortURL:  tt.fields.origToShortURL,
				usersURLS:       tt.fields.usersURLS,
				batchBuffer:     tt.fields.batchBuffer,
				batchCounter:    tt.fields.batchCounter,
				batchSize:       tt.fields.batchSize,
				storageFilePath: tt.fields.storageFilePath,
			}
			tt.wantErr(t, m.StoreURL(tt.args.ctx, tt.args.originalURL, tt.args.shortURL), fmt.Sprintf("StoreURL(%v, %v)", tt.args.originalURL, tt.args.shortURL))
			defer os.Remove(tt.fields.storageFilePath)
		})
	}
}

func createTempFilePath(t *testing.T) string {
	tempFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	_, err = tempFile.Write([]byte(`{"user_id":"e774844b-5895-4b08-b867-50480263f75b","short_url":"short1","original_url":"original1"}`))
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempPath := tempFile.Name()
	tempFile.Close()
	return tempPath
}

func TestURLInMemoryRepo_StoreBatchURL(t *testing.T) {
	tests := []struct {
		name             string
		batchURLtoStores map[string]string
		expectedError    error
	}{
		{
			name: "Valid batch of URLs",
			batchURLtoStores: map[string]string{
				"short1": "http://example1.com",
				"short2": "http://example2.com",
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mocking the context for testing
			ctx := context.WithValue(context.Background(), models.UserIDKey, uuid.New())
			// Создаем временный файл для теста
			tempFile, err := os.CreateTemp("", "temp_test_file.json")
			assert.NoError(t, err)
			defer func() {
				tempFile.Close()
				// Удаляем временный файл после теста
				err := os.Remove(tempFile.Name())
				assert.NoError(t, err)
			}()

			// Создаем репозиторий
			repo := NewURLInMemoryRepo(tempFile.Name())
			// Call the method under test
			err = repo.StoreBatchURL(ctx, tt.batchURLtoStores)

			// Check the result
			if (err != nil && tt.expectedError == nil) || (err == nil && tt.expectedError != nil) || (err != nil && tt.expectedError != nil && err.Error() != tt.expectedError.Error()) {
				t.Errorf("Unexpected result. Expected error: %v, Got error: %v", tt.expectedError, err)
			}

			// Optionally, you can check the internal state of m after the operation.
			// For example, check if the URLs are stored correctly in m.usersURLS and m.shortToOrigURL.
		})
	}
}

func TestURLInMemoryRepo_GetShortBatchURL(t *testing.T) {
	tests := []struct {
		name              string
		batchURLRequests  []models.URLRequest
		expectedShortURLs map[string]string
		expectedErr       bool
	}{
		{
			name: "Batch URLs found",
			batchURLRequests: []models.URLRequest{
				{OriginalURL: "http://example1.com"},
				{OriginalURL: "http://example2.com"},
			},
			expectedShortURLs: map[string]string{
				"http://example1.com": "short1",
				"http://example2.com": "short2",
			},
			expectedErr: false,
		},
	}

	for _, tt := range tests {
		createTempFilePath(t)
		t.Run(tt.name, func(t *testing.T) {
			// Создаем временный файл для теста
			tempFile, err := os.CreateTemp("", "temp_test_file.json")
			assert.NoError(t, err)
			defer func() {
				tempFile.Close()
				// Удаляем временный файл после теста
				err := os.Remove(tempFile.Name())
				assert.NoError(t, err)
			}()

			// Создаем репозиторий
			repo := NewURLInMemoryRepo(tempFile.Name())
			defer repo.SaveBatchToFile() // Сохраняем оставшиеся данные перед завершением теста

			// Добавляем тестовые URL в репозиторий
			for i, req := range tt.batchURLRequests {
				shortURL := fmt.Sprintf("short%d", i+1)
				err := repo.StoreURL(context.Background(), req.OriginalURL, shortURL)
				assert.NoError(t, err)
			}

			// Вызываем метод, который мы тестируем
			shortURLs, err := repo.GetShortBatchURL(context.Background(), tt.batchURLRequests)

			// Проверяем ошибку
			assert.Equal(t, tt.expectedErr, err != nil)

			// Проверяем, что полученные короткие URL соответствуют ожидаемым
			assert.Equal(t, tt.expectedShortURLs, shortURLs)

		})
	}
}

func TestURLInMemoryRepo_GetUserURLS(t *testing.T) {
	tests := []struct {
		name        string
		userID      uuid.UUID
		expectedErr bool
		expectedURL *models.URL
	}{
		{
			name:        "User has URLs",
			userID:      uuid.New(),
			expectedErr: false,
			expectedURL: &models.URL{
				ShortURL:    "http://short.com",
				OriginalURL: "http://example.com",
			},
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем временный файл для теста
			tempFile, err := os.CreateTemp("", "temp_test_file.json")
			assert.NoError(t, err)
			defer func() {
				tempFile.Close()
				// Удаляем временный файл после теста
				err := os.Remove(tempFile.Name())
				assert.NoError(t, err)
			}()

			// Создаем репозиторий
			repo := NewURLInMemoryRepo(tempFile.Name())
			defer repo.SaveBatchToFile() // Сохраняем оставшиеся данные перед завершением теста

			// Создаем контекст с указанным userID
			ctx := context.WithValue(context.Background(), models.UserIDKey, tt.userID)

			// Добавляем тестовый URL для пользователя
			err = repo.StoreURL(ctx, "http://example.com", "http://short.com")
			assert.NoError(t, err)

			// Вызываем метод, который мы тестируем
			userURLs, err := repo.GetUserURLS(ctx)

			// Проверяем ошибку
			assert.Equal(t, tt.expectedErr, err != nil)

			// Проверяем, что полученные URL соответствуют ожидаемым
			assert.Equal(t, len(userURLs), 1)

			// Сравниваем URL по полям
			assert.Equal(t, userURLs[0].ShortURL, tt.expectedURL.ShortURL)
			assert.Equal(t, userURLs[0].OriginalURL, tt.expectedURL.OriginalURL)
		})
	}
}
