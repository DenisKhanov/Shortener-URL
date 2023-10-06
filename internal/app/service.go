package app

import "strings"

// GetShortURL returns the short URL ("http://localhost:8080/"+shortURL)
func GetShortURL(url string) string {
	value, exists := store.origToShortURL[url]
	if exists {
		return "http://localhost:8080/" + value
	} else {
		store.unicalID++
		shortURL := base62Encode(store.unicalID)
		store.shortToOrigURL[shortURL] = url
		store.origToShortURL[url] = shortURL
		return "http://localhost:8080/" + shortURL
	}
}

// GetOriginURL returns the origin URL for the given short URL
func GetOriginURL(shortURL string) (string, bool) {
	originURL, exists := store.shortToOrigURL[shortURL]
	return originURL, exists

}

// base62Encode returns the base 64 encoded string
func base62Encode(n int) string {
	chars := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	if n == 0 {
		return string(chars[0])
	}
	var shortURL strings.Builder
	for n > 0 {
		shortURL.WriteString(string(chars[n%62]))
		n = n / 62
	}
	return shortURL.String()
}
