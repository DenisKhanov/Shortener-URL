package main

import (
	"fmt"
	"github.com/DenisKhanov/shorterURL/internal/app/config"
	"github.com/DenisKhanov/shorterURL/internal/app/handlers"
	"github.com/DenisKhanov/shorterURL/internal/app/repositoryes"
	"github.com/DenisKhanov/shorterURL/internal/app/services"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	cfg := config.NewConfig()
	fmt.Println("Server Address:", cfg.EnvServAdr)
	fmt.Println("Base URL:", cfg.EnvBaseURL)
	repository := repositoryes.NewRepository(make(map[string]string), make(map[string]string))

	myService := services.NewServices(repository, services.Services{}, cfg.EnvBaseURL)
	myHandler := handlers.NewHandlers(myService)

	r := mux.NewRouter()
	r.HandleFunc("/", myHandler.PostURL)
	r.HandleFunc("/{id}", myHandler.GetURL).Methods("GET")

	fmt.Println("Server started on", cfg.EnvServAdr)

	if err := http.ListenAndServe(cfg.EnvServAdr, r); err != nil {
		panic(err)
	}
}
