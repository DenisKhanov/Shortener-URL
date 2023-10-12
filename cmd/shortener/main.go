package main

import (
	"fmt"
	"github.com/DenisKhanov/shorterURL/internal/app/config"
	"github.com/DenisKhanov/shorterURL/internal/app/handlers"
	"github.com/DenisKhanov/shorterURL/internal/app/service"
	"github.com/DenisKhanov/shorterURL/internal/app/storage"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	cfg := config.NewConfig()
	fmt.Println("Server Address:", cfg.EnvServAdr)
	fmt.Println("Base URL:", cfg.EnvBaseURL)
	repository := storage.NewRepository(214134121, make(map[string]string), make(map[string]string))
	var s service.Service
	var myService = service.NewServices(repository, s, cfg.EnvBaseURL)
	var handler handlers.Handler
	var myHandler = handlers.NewHandlers(*myService, handler)

	r := mux.NewRouter()
	r.HandleFunc("/", myHandler.PostURL)
	r.HandleFunc("/{id}", myHandler.GetURL).Methods("GET")

	fmt.Println("Server started on", cfg.EnvServAdr)

	if err := http.ListenAndServe(cfg.EnvServAdr, r); err != nil {
		panic(err)
	}
}
