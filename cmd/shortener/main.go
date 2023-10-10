package main

import (
	"github.com/DenisKhanov/shorterURL/internal/app/handlers"
	"github.com/DenisKhanov/shorterURL/internal/app/service"
	"github.com/DenisKhanov/shorterURL/internal/app/storage"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	repository := storage.NewRepository(214134121, make(map[string]string), make(map[string]string))
	var s service.Service
	var myService = service.NewServices(repository, s)
	var handler handlers.Handler
	var myHandler = handlers.NewHandlers(*myService, handler)

	r := mux.NewRouter()
	r.HandleFunc("/", myHandler.PostURL)
	r.HandleFunc("/{id}", myHandler.GetURL).Methods("GET")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)

	}
}
