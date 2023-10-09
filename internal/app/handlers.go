package main

import (
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

//go:generate mockgen -source=handlers.go -destination=mocks/handlers_mock.go -package=mocks
type Handler interface {
	URL(w http.ResponseWriter, r *http.Request)
}
type Handlers struct {
	service Service
	handler Handler
}

func NewHandlers(service Service, handler Handler) *Handlers {
	return &Handlers{
		service: service,
		handler: handler,
	}
}
func (h Handlers) PostURL(w http.ResponseWriter, r *http.Request) {
	//if r.Method == http.MethodPost {
	url, _ := io.ReadAll(r.Body)
	r.Body.Close()
	if len(url) == 0 {
		w.WriteHeader(http.StatusBadRequest)

	} else {
		shortURL := h.service.GetShortURL(string(url))
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(shortURL))

	}
	//} else {
	//	w.WriteHeader(http.StatusBadRequest)
	//}

}
func (h Handlers) GetURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortURL := vars["id"]
	originURL, err := h.service.GetOriginURL(shortURL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.Header().Set("Location", originURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
