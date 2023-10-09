package app

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRepository_StoreURL_GetOriginalURL_GetShortURL_GetID(t *testing.T) {
	repo := NewRepository(0, make(map[string]string), make(map[string]string))
	originalURL := "https://example.com"
	shortURL := "abc123"

	repo.StoreURL(originalURL, shortURL)
	retrievedOriginalURL, exists := repo.GetOriginalURL(shortURL)
	assert.True(t, exists)
	assert.Equal(t, originalURL, retrievedOriginalURL)

	retrievedShortURL, exists := repo.GetShortURL(originalURL)
	assert.True(t, exists)
	assert.Equal(t, shortURL, retrievedShortURL)

	assert.Equal(t, 1, repo.GetID())
}
func TestNewDumpURL(t *testing.T) {
	id := 1
	shortToOrigURL := map[string]string{"short": "original"}
	origToShortURL := map[string]string{"original": "short"}

	result := NewDumpURL(id, shortToOrigURL, origToShortURL)

	assert.Equal(t, id, result.id)
	assert.Equal(t, shortToOrigURL, result.shortToOrigURL)
	assert.Equal(t, origToShortURL, result.origToShortURL)
}
