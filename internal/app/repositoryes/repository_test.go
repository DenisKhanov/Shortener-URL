package repositoryes

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewRepository(t *testing.T) {
	tests := []struct {
		name string
		want *URLInMemoryRepo
	}{
		{
			name: "Valid args",
			want: &URLInMemoryRepo{
				shortToOrigURL:  map[string]string{"short1": "original1"},
				origToShortURL:  map[string]string{"original1": "short1"},
				lastUUID:        1,
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
	}
	type args struct {
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
			name:    "valid get original URL",
			fields:  fields{map[string]string{"short": "original"}, map[string]string{"original": "short"}},
			args:    args{shortURL: "short"},
			want:    "original",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &URLInMemoryRepo{
				shortToOrigURL: tt.fields.shortToOrigURL,
				origToShortURL: tt.fields.origToShortURL,
			}
			got, err := d.GetOriginalURLFromDB(tt.args.shortURL)
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
	}
	type args struct {
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
			name:    "valid get short URL",
			fields:  fields{map[string]string{"short": "original"}, map[string]string{"original": "short"}},
			args:    args{originalURL: "original"},
			want:    "short",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &URLInMemoryRepo{
				shortToOrigURL: tt.fields.shortToOrigURL,
				origToShortURL: tt.fields.origToShortURL,
			}
			got, err := d.GetShortURLFromDB(tt.args.originalURL)
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
		lastUUID        int
		batchBuffer     []URLInFileRepo
		batchCounter    int
		batchSize       int
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
				shortToOrigURL:  map[string]string{"short1": "original1"},
				origToShortURL:  map[string]string{"original1": "short1"},
				lastUUID:        1,
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
				lastUUID:        tt.fields.lastUUID,
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
		lastUUID        int
		batchBuffer     []URLInFileRepo
		batchCounter    int
		batchSize       int
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
				lastUUID:       1,
				batchBuffer: []URLInFileRepo{
					{
						UUID:        "1",
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
				lastUUID:        tt.fields.lastUUID,
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
		lastUUID        int
		batchBuffer     []URLInFileRepo
		batchCounter    int
		batchSize       int
		storageFilePath string
	}
	type args struct {
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
				lastUUID:       1,
				batchBuffer: []URLInFileRepo{
					{
						UUID:        "2",
						ShortURL:    "short2",
						OriginalURL: "original2",
					},
				},
				batchCounter:    1,
				batchSize:       100,
				storageFilePath: createTempFilePath(t),
			},
			args:    args{originalURL: "original2", shortURL: "short2"},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &URLInMemoryRepo{
				shortToOrigURL:  tt.fields.shortToOrigURL,
				origToShortURL:  tt.fields.origToShortURL,
				lastUUID:        tt.fields.lastUUID,
				batchBuffer:     tt.fields.batchBuffer,
				batchCounter:    tt.fields.batchCounter,
				batchSize:       tt.fields.batchSize,
				storageFilePath: tt.fields.storageFilePath,
			}
			tt.wantErr(t, m.StoreURLInDB(tt.args.originalURL, tt.args.shortURL), fmt.Sprintf("StoreURLInDB(%v, %v)", tt.args.originalURL, tt.args.shortURL))
			defer os.Remove(tt.fields.storageFilePath)
		})
	}
}

func createTempFilePath(t *testing.T) string {
	tempFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	_, err = tempFile.Write([]byte(`{"uuid":"1","short_url":"short1","original_url":"original1"}`))
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempPath := tempFile.Name()
	tempFile.Close()
	return tempPath
}
