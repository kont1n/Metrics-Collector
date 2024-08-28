package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"Metrics-Collector/internal/api"
	"Metrics-Collector/internal/storage"
)

const (
	serverPort = 8080
)

var (
	err error
)

func main() {
	store := storage.NewMemStorage()

	router := chi.NewRouter()
	router.Post("/update/{type}/{metric}/{value}", api.PostMetric(store))
	router.Get("/value/{type}/{metric}", api.GetMetrics(store))
	//router.Handle("/*", http.FileServer(http.Dir("web")))

	fmt.Println("Server started on localhost port", serverPort)
	if err = http.ListenAndServe(fmt.Sprintf(":%d", serverPort), router); err != nil {
		fmt.Println("Web server error:", err.Error())
		return
	}
}
