package main

import (
	"fmt"
	"net/http"

	"Metrics-Collector/internal/config"
	"github.com/go-chi/chi/v5"

	"Metrics-Collector/internal/api"
	"Metrics-Collector/internal/storage"
)

var (
	err error
)

func main() {
	host := config.ParseServerConfig()

	store := storage.NewMemStorage()

	router := chi.NewRouter()
	router.Post("/update/{type}/{metric}/{value}", api.PostMetric(store))
	router.Get("/value/{type}/{metric}", api.GetMetrics(store))
	router.Get("/", api.IndexHandler(store))

	fmt.Printf("Server started on %s\n", host)
	if err = http.ListenAndServe(host, router); err != nil {
		fmt.Println("Web server error:", err.Error())
		return
	}
}
