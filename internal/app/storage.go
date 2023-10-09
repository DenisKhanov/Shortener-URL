package main

//go:generate mockgen -source=storage.go -destination=mocks/storage_mock.go -package=mocks

type Repository interface {
	StoreURL(originalURL, shortURL string)
	GetShortURL(originalURL string) (string, bool)
	GetOriginalURL(shortURL string) (string, bool)
	GetID() int
}

type repositoryURL struct {
	id             int
	shortToOrigURL map[string]string
	origToShortURL map[string]string
}

func NewRepository(id int, shortToOrigURL map[string]string, origToShortURL map[string]string) Repository {
	storage := repositoryURL{
		id:             id,
		shortToOrigURL: shortToOrigURL,
		origToShortURL: origToShortURL,
	}
	return &storage
}

func NewDumpURL(id int, shortToOrigURL, origToShortURL map[string]string) *repositoryURL {
	return &repositoryURL{
		id:             id,
		shortToOrigURL: shortToOrigURL,
		origToShortURL: origToShortURL,
	}
}
func (d *repositoryURL) StoreURL(originalURL, shortURL string) {
	d.origToShortURL[originalURL] = shortURL
	d.shortToOrigURL[shortURL] = originalURL
}
func (d *repositoryURL) GetOriginalURL(shortURL string) (string, bool) {
	originalURL, exists := d.shortToOrigURL[shortURL]
	return originalURL, exists
}
func (d *repositoryURL) GetShortURL(originalURL string) (string, bool) {
	shortURL, exists := d.origToShortURL[originalURL]
	return shortURL, exists
}
func (d *repositoryURL) GetID() int {
	d.id++
	return d.id
}
