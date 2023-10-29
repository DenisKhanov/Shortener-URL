package repositoryes

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewURLInFileRepo(t *testing.T) {
	type args struct {
		storageFilePath string
	}
	tests := []struct {
		name string
		args args
		want *URLInFileRepo
	}{
		{
			name: "Valid args",
			args: args{storageFilePath: "/tmp/test.json"},
			want: &URLInFileRepo{"", "", "", "/tmp/test.json"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewURLInFileRepo(tt.args.storageFilePath), "NewURLInFileRepo(%v)", tt.args.storageFilePath)
		})
	}
}

func TestURLInFileRepo_GetOriginalURLFromDB(t *testing.T) {
	type fields struct {
		UUID            string
		ShortURL        string
		OriginalURL     string
		storageFilePath string
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
			fields:  fields{"1", "short", "original", createTempFilePath(t)},
			args:    args{shortURL: "short"},
			want:    "original",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &URLInFileRepo{
				UUID:            tt.fields.UUID,
				ShortURL:        tt.fields.ShortURL,
				OriginalURL:     tt.fields.OriginalURL,
				storageFilePath: tt.fields.storageFilePath,
			}

			got, err := r.GetOriginalURLFromDB(tt.args.shortURL)
			if !tt.wantErr(t, err, fmt.Sprintf("GetOriginalURLFromDB(%v)", tt.args.shortURL)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetOriginalURLFromDB(%v)", tt.args.shortURL)
			defer os.Remove(tt.fields.storageFilePath)
		})
	}
}
func TestURLInFileRepo_StoreURLSInDB(t *testing.T) {
	type fields struct {
		UUID            string
		ShortURL        string
		OriginalURL     string
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
			name:    "valid store short URL",
			fields:  fields{"1", "short", "original", createTempFilePath(t)},
			args:    args{originalURL: "original", shortURL: "short"},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &URLInFileRepo{
				UUID:            tt.fields.UUID,
				ShortURL:        tt.fields.ShortURL,
				OriginalURL:     tt.fields.OriginalURL,
				storageFilePath: tt.fields.storageFilePath,
			}
			tt.wantErr(t, r.StoreURLSInDB(tt.args.originalURL, tt.args.shortURL), fmt.Sprintf("StoreURLSInDB(%v, %v)", tt.args.originalURL, tt.args.shortURL))
		})
	}
}

func TestURLInFileRepo_GetShortURLFromDB(t *testing.T) {
	type fields struct {
		UUID            string
		ShortURL        string
		OriginalURL     string
		storageFilePath string
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
			fields:  fields{"1", "short", "original", createTempFilePath(t)},
			args:    args{originalURL: "original"},
			want:    "short",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &URLInFileRepo{
				UUID:            tt.fields.UUID,
				ShortURL:        tt.fields.ShortURL,
				OriginalURL:     tt.fields.OriginalURL,
				storageFilePath: tt.fields.storageFilePath,
			}
			got, err := r.GetShortURLFromDB(tt.args.originalURL)
			if !tt.wantErr(t, err, fmt.Sprintf("GetShortURLFromDB(%v)", tt.args.originalURL)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetShortURLFromDB(%v)", tt.args.originalURL)
			defer os.Remove(tt.fields.storageFilePath)
		})
	}
}

func TestURLInFileRepo_LoadLastUUID(t *testing.T) {
	type fields struct {
		UUID            string
		ShortURL        string
		OriginalURL     string
		storageFilePath string
	}
	tests := []struct {
		name     string
		fields   fields
		expected int
	}{
		{
			name:     "valid last uuid",
			fields:   fields{"1", "short", "original", createTempFilePath(t)},
			expected: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &URLInFileRepo{
				UUID:            tt.fields.UUID,
				ShortURL:        tt.fields.ShortURL,
				OriginalURL:     tt.fields.OriginalURL,
				storageFilePath: tt.fields.storageFilePath,
			}
			assert.Equal(t, tt.expected, r.LoadLastUUID())
			defer os.Remove(tt.fields.storageFilePath)
		})
	}
}

func Test_writeLine(t *testing.T) {
	type args struct {
		file *os.File
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		wantN   int
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "valid write line",
			args: args{
				file: createTempFile(t),
				data: []byte(`{"uuid":"1","short_url":"lqPaDlz","original_url":"https://test.ru"}`),
			},
			wantN:   len(`{"uuid":"1","short_url":"lqPaDlz","original_url":"https://test.ru"}`) + 1,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotN, err := writeLine(tt.args.file, tt.args.data)
			if !tt.wantErr(t, err, fmt.Sprintf("writeLine(%v, %v)", tt.args.file, tt.args.data)) {
				return
			}
			assert.Equalf(t, tt.wantN, gotN, "writeLine(%v, %v)", tt.args.file, tt.args.data)
			cleanupTempFile(t, tt.args.file)
		})
	}
}
func createTempFile(t *testing.T) *os.File {
	tempFile, err := os.CreateTemp("", "test.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	return tempFile
}
func createTempFilePath(t *testing.T) string {
	tempFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	_, err = tempFile.Write([]byte(`{"uuid":"1","short_url":"short","original_url":"original"}`))
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempPath := tempFile.Name()
	tempFile.Close()
	return tempPath
}

func cleanupTempFile(t *testing.T, file *os.File) {
	file.Close()
	os.Remove(file.Name())
}
