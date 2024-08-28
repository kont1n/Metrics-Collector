package api

import (
	"fmt"
	"html"
	"net/http"
	"sort"
	"strconv"

	"github.com/go-chi/chi/v5"

	"Metrics-Collector/internal/storage"
)

func PostMetric(store *storage.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "type")
		metricName := chi.URLParam(r, "metric")
		metricValue := chi.URLParam(r, "value")

		fmt.Println("handler:", metricType, metricName, metricValue)

		if metricName == "" {
			w.Header().Set("Content-Type", "text/plain")
			http.Error(w, "incorrect metric name", http.StatusNotFound)
			return
		}

		if metricValue == "" {
			w.Header().Set("Content-Type", "text/plain")
			http.Error(w, "incorrect metric value", http.StatusBadRequest)
			return
		}

		switch metricType {
		case "gauge":
			{
				value, err := strconv.ParseFloat(metricValue, 64)
				if err != nil {
					w.Header().Set("Content-Type", "text/plain")
					http.Error(w, "incorrect metric value", http.StatusBadRequest)
					return
				}
				store.SetGauge(metricName, value)
				w.Header().Set("Content-Type", "text/plain")
				w.WriteHeader(http.StatusOK)
			}
		case "counter":
			{
				value, err := strconv.ParseInt(metricValue, 10, 64)
				if err != nil {
					w.Header().Set("Content-Type", "text/plain")
					http.Error(w, "incorrect metric value", http.StatusBadRequest)
					return
				}
				store.SetCounter(metricName, value)
				w.Header().Set("Content-Type", "text/plain")
				w.WriteHeader(http.StatusOK)
			}
		default:
			{
				w.Header().Set("Content-Type", "text/plain")
				http.Error(w, "incorrect metric type", http.StatusBadRequest)
			}
		}
	}
}

func GetMetrics(store *storage.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var answer string
		w.Header().Set("Content-Type", "text/plain")

		metricType := chi.URLParam(r, "type")
		metricName := chi.URLParam(r, "metric")

		switch metricType {
		case "gauge":
			{
				value, exists := store.GetGauge(metricName)
				if !exists {
					http.Error(w, "unknown metric", http.StatusNotFound)
					return
				}
				answer = fmt.Sprintf("%.3f", value)
			}
		case "counter":
			{
				value, exists := store.GetCounter(metricName)
				if !exists {
					http.Error(w, "unknown metric", http.StatusNotFound)
					return
				}
				answer = fmt.Sprintf("%d", value)
			}
		default:
			{
				w.Header().Set("Content-Type", "text/plain")
				http.Error(w, "incorrect metric type", http.StatusBadRequest)
			}
		}
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(answer))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func IndexHandler(store *storage.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var result string

		for metric, value := range store.GetCounters() {
			result += fmt.Sprintf("%s: %d\n", metric, value)
		}

		m := store.GetGauges()

		keys := make([]string, 0, len(m))
		for k := range store.GetGauges() {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			result += fmt.Sprintf("%s: %.3f\n", k, m[k])
		}

		htmlString := "<html><head><title>Metrics</title></head><body><pre>" +
			html.EscapeString(result) +
			"</pre></body></html>"

		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)

		if _, err := w.Write([]byte(htmlString)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
