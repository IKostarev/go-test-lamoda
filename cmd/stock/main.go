package main

import (
	_ "github.com/lib/pq"
	"go-test-lamoda/internal/database"
	"go-test-lamoda/internal/handlers"
	"log"
	"net/http"
)

func main() {
	_, err := database.NewPostgresDB(database.Config{
		Host:     "localhost",
		Port:     "5432",
		Username: "lamoda_test",
		Password: "test",
		DBName:   "test",
		SSLMode:  "disable",
	})
	if err != nil {
		log.Fatalf("Ошибка при инициализации базы данных: %s", err.Error())
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.HelloHandler)
	mux.HandleFunc("/reserve", handlers.ReserveToStockHandler)
	mux.HandleFunc("/release", handlers.ReleaseHandler)
	mux.HandleFunc("/count", handlers.GetCountHandler)

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
