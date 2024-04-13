package repositories

import (
	"context"
	"fmt"
	"github.com/DenisKhanov/shorterURL/internal/app/models"
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

func TestRepositoryURL_GetOriginalURLFromDB(t *testing.T) {
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
			got, err := d.GetOriginalURLFromDB(tt.args.ctx, tt.args.shortURL)
			if !tt.wantErr(t, err, fmt.Sprintf("GetOriginalURLFromDB(%v)", tt.args.shortURL)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetOriginalURLFromDB(%v)", tt.args.shortURL)
		})
	}
}

func TestRepositoryURL_GetShortURLFromDB(t *testing.T) {
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
			got, err := d.GetShortURLFromDB(tt.args.ctx, tt.args.originalURL)
			if !tt.wantErr(t, err, fmt.Sprintf("GetShortURLFromDB(%v)", tt.args.originalURL)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetShortURLFromDB(%v)", tt.args.originalURL)
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

func TestURLInMemoryRepo_StoreURLSInDB(t *testing.T) {
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
			tt.wantErr(t, m.StoreURLInDB(tt.args.ctx, tt.args.originalURL, tt.args.shortURL), fmt.Sprintf("StoreURLInDB(%v, %v)", tt.args.originalURL, tt.args.shortURL))
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
