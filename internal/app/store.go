package app

type URLStore struct {
	unicalID       int
	shortToOrigURL map[string]string
	origToShortURL map[string]string
}

var store = NewDumpURL(1285663434, make(map[string]string), make(map[string]string))

func NewDumpURL(unicalID int, shortToOrigURL, origToShortURL map[string]string) *URLStore {
	return &URLStore{
		unicalID:       unicalID,
		shortToOrigURL: shortToOrigURL,
		origToShortURL: origToShortURL,
	}
}
