package main

import (
	"github.com/DenisKhanov/shorterURL/internal/app"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	storage := app.NewRepository(214134121, make(map[string]string), make(map[string]string))
	var service app.Service
	var myService = app.NewServices(storage, service)
	var handler app.Handler
	var myHandler = app.NewHandlers(*myService, handler)

	r := mux.NewRouter()
	r.HandleFunc("/", myHandler.PostURL)
	r.HandleFunc("/{id}", myHandler.GetURL).Methods("GET")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)

	}
}
