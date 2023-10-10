package storage

import "errors"

//go:generate mockgen -source=storage.go -destination=mocks/storage_mock.go -package=mocks

type Repository interface {
	StoreURL(originalURL, shortURL string) error
	GetShortURL(originalURL string) (string, bool)
	GetOriginalURL(shortURL string) (string, bool)
	GetID() int
}

type RepositoryURL struct {
	id             int
	shortToOrigURL map[string]string
	origToShortURL map[string]string
}

func NewRepository(id int, shortToOrigURL map[string]string, origToShortURL map[string]string) Repository {
	storage := RepositoryURL{
		id:             id,
		shortToOrigURL: shortToOrigURL,
		origToShortURL: origToShortURL,
	}
	return &storage
}

func NewDumpURL(id int, shortToOrigURL, origToShortURL map[string]string) *RepositoryURL {
	return &RepositoryURL{
		id:             id,
		shortToOrigURL: shortToOrigURL,
		origToShortURL: origToShortURL,
	}
}
func (d *RepositoryURL) StoreURL(originalURL, shortURL string) error {
	d.origToShortURL[originalURL] = shortURL
	d.shortToOrigURL[shortURL] = originalURL
	if d.origToShortURL[originalURL] == "" || d.shortToOrigURL[shortURL] == "" {
		return errors.New("error saving shortUrl")
	}
	return nil
}
func (d *RepositoryURL) GetOriginalURL(shortURL string) (string, bool) {
	originalURL, exists := d.shortToOrigURL[shortURL]
	return originalURL, exists
}
func (d *RepositoryURL) GetShortURL(originalURL string) (string, bool) {
	shortURL, exists := d.origToShortURL[originalURL]
	return shortURL, exists
}
func (d *RepositoryURL) GetID() int {
	d.id++
	return d.id
}
