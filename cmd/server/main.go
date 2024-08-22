package main

import (
	"Metrics-Collector/internal/api"
	"Metrics-Collector/internal/storage"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

const (
	serverPort = 8080
)

var (
	err error
)

func main() {
	store := storage.NewMemStorage()

	r := chi.NewRouter()
	r.Post("/update/{type}/{metric}/{value}", api.PostMetric(*store))

	fmt.Println("Server started on localhost port", serverPort)
	if err = http.ListenAndServe(fmt.Sprintf(":%d", serverPort), nil); err != nil {
		fmt.Println("Web server error:", err.Error())
		return
	}
}
