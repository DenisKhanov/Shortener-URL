package repositoryes

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRepository(t *testing.T) {
	id := 214134121
	shortToOrigURL := map[string]string{"short": "original"}
	origToShortURL := map[string]string{"original": "short"}

	repository := NewRepository(id, shortToOrigURL, origToShortURL)

	assert.Equal(t, id, repository.id)
	assert.Equal(t, shortToOrigURL, repository.shortToOrigURL)
	assert.Equal(t, origToShortURL, repository.origToShortURL)
}

func TestRepository(t *testing.T) {
	repository := RepositoryURL{id: 214134121, shortToOrigURL: make(map[string]string), origToShortURL: make(map[string]string)}
	tests := []struct {
		name        string
		originalURL string
		shortURL    string
		exists      bool
		validID     int
		repoMethod  RepositoryURL
	}{
		{
			name:        "All methods is valid",
			originalURL: "http://original.url",
			shortURL:    "A4UUE",
			exists:      true,
			validID:     214134122,
			repoMethod:  repository,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository = RepositoryURL{id: 214134121, shortToOrigURL: make(map[string]string), origToShortURL: make(map[string]string)}

			assert.NoError(t, tt.repoMethod.StoreURLSInDB(tt.originalURL, tt.shortURL))

			resultOriginal, existsOriginal := tt.repoMethod.GetOriginalURLFromDB(tt.shortURL)
			if assert.Equal(t, tt.exists, existsOriginal) {
				assert.Equal(t, tt.originalURL, resultOriginal)
			}
			resultShort, existsShort := tt.repoMethod.GetShortURLFromDB(tt.originalURL)
			if assert.Equal(t, tt.exists, existsShort) {
				assert.Equal(t, tt.shortURL, resultShort)
			}
			assert.Equal(t, tt.validID, tt.repoMethod.GetIDFromDB())
		})
	}
}
