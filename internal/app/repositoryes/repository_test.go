package repositoryes

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRepository(t *testing.T) {
	tests := []struct {
		name string
		want *URLInMemoryRepo
	}{
		{
			name: "Valid args",
			want: &URLInMemoryRepo{map[string]string{}, map[string]string{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewURLInMemoryRepo(), "NewURLInMemoryRepo()")
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

func TestRepositoryURL_StoreURLSInDB(t *testing.T) {
	type fields struct {
		shortToOrigURL map[string]string
		origToShortURL map[string]string
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
			fields:  fields{map[string]string{"short": "original"}, map[string]string{"original": "short"}},
			args:    args{originalURL: "original", shortURL: "short"},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &URLInMemoryRepo{
				shortToOrigURL: tt.fields.shortToOrigURL,
				origToShortURL: tt.fields.origToShortURL,
			}
			tt.wantErr(t, d.StoreURLSInDB(tt.args.originalURL, tt.args.shortURL), fmt.Sprintf("StoreURLSInDB(%v, %v)", tt.args.originalURL, tt.args.shortURL))
		})
	}
}
