package repositoryes

import "errors"

type RepositoryURL struct {
	id             int
	shortToOrigURL map[string]string
	origToShortURL map[string]string
}

func NewRepository(id int, shortToOrigURL map[string]string, origToShortURL map[string]string) *RepositoryURL {
	storage := RepositoryURL{
		id:             id,
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
func (d *RepositoryURL) GetOriginalURLFromDB(shortURL string) (string, bool) {
	originalURL, exists := d.shortToOrigURL[shortURL]
	return originalURL, exists
}
func (d *RepositoryURL) GetShortURLFromDB(originalURL string) (string, bool) {
	shortURL, exists := d.origToShortURL[originalURL]
	return shortURL, exists
}
func (d *RepositoryURL) GetIDFromDB() int {
	d.id++
	return d.id
}
