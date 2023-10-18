package repositoryes

import (
	"errors"
)

type RepositoryURL struct {
	shortToOrigURL map[string]string
	origToShortURL map[string]string
}

func NewRepository(shortToOrigURL map[string]string, origToShortURL map[string]string) *RepositoryURL {
	storage := RepositoryURL{
		shortToOrigURL: shortToOrigURL,
		origToShortURL: origToShortURL,
	}
	return &storage
}

func (d *RepositoryURL) StoreURLSInDB(originalURL, shortURL string) error {
	d.origToShortURL[originalURL] = shortURL
	d.shortToOrigURL[shortURL] = originalURL
	if d.origToShortURL[originalURL] == "" || d.shortToOrigURL[shortURL] == "" {
		return errors.New("error saving shortUrl")
	}
	return nil
}
func (d *RepositoryURL) GetOriginalURLFromDB(shortURL string) (string, error) {
	originalURL, exists := d.shortToOrigURL[shortURL]
	if !exists {
		return "", errors.New("original URL not found")
	}
	return originalURL, nil
}
func (d *RepositoryURL) GetShortURLFromDB(originalURL string) (string, error) {
	shortURL, exists := d.origToShortURL[originalURL]
	if !exists {
		return "", errors.New("short URL not found")
	}
	return shortURL, nil
}
