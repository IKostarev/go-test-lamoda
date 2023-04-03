package main

import (
	"go-test-lamoda/internal/handlers"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.HelloHandler)
	mux.HandleFunc("/reserve", handlers.ReserveToStockHandler)
	mux.HandleFunc("/release", handlers.ReleaseHandler)
	mux.HandleFunc("/count", handlers.GetCountHandler)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
