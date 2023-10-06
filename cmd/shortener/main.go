package main

import (
	"github.com/DenisKhanov/shorterURL/internal/app"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.URL)
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}

}
