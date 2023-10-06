package app

import (
	"io"
	"net/http"
)

func URL(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		url, _ := io.ReadAll(r.Body)
		if len(url) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Ссылка в теле не передана"))

		} else {
			shortURL := GetShortURL(string(url))
			w.WriteHeader(http.StatusCreated)
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte(shortURL))

		}
		if r.Method == http.MethodGet {
			path := r.URL.Path[1:]
			originURL, exists := GetOriginURL(path)
			if !exists {
				http.NotFound(w, r)
			} else {
				w.Header().Set("Location", originURL)
				w.WriteHeader(http.StatusTemporaryRedirect)
			}
		}

	}
}
