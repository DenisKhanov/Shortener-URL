package repositoryes

import "errors"

type URLInMemoryRepo struct {
	shortToOrigURL map[string]string
	origToShortURL map[string]string
}

func NewURLInMemoryRepo() *URLInMemoryRepo {
	storage := URLInMemoryRepo{
		shortToOrigURL: make(map[string]string),
		origToShortURL: make(map[string]string),
	}
	return &storage
}

func (d *URLInMemoryRepo) StoreURLSInDB(originalURL, shortURL string) error {
	d.origToShortURL[originalURL] = shortURL
	d.shortToOrigURL[shortURL] = originalURL
	if d.origToShortURL[originalURL] == "" || d.shortToOrigURL[shortURL] == "" {
		return errors.New("error saving shortUrl")
	}
	return nil
}
func (d *URLInMemoryRepo) GetOriginalURLFromDB(shortURL string) (string, error) {
	originalURL, exists := d.shortToOrigURL[shortURL]
	if !exists {
		return "", errors.New("original URL not found")
	}
	return originalURL, nil
}
func (d *URLInMemoryRepo) GetShortURLFromDB(originalURL string) (string, error) {
	shortURL, exists := d.origToShortURL[originalURL]
	if !exists {
		return "", errors.New("short URL not found")
	}
	return shortURL, nil
}
