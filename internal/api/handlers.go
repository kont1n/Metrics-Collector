package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"sort"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"Metrics-Collector/internal/models"
	"Metrics-Collector/internal/service"
)

const (
	TextHTML        = "text/html"
	TextPlain       = "text/plain"
	ApplicationJSON = "application/json"
)

type Handler struct {
	service *service.Service
	loger   *zap.SugaredLogger
	Handlers
}

type Handlers interface {
	postMetric(w http.ResponseWriter, r *http.Request)
	getMetric(w http.ResponseWriter, r *http.Request)
	indexHandler(w http.ResponseWriter, r *http.Request)
	postJSONMetric(w http.ResponseWriter, r *http.Request)
	getJSONMetric(w http.ResponseWriter, r *http.Request)
	withJSON(w http.ResponseWriter, v any, status int, reqID string)
	jsonError(w http.ResponseWriter, error string, code int, reqID string)
}

func NewHandler(service *service.Service, loger *zap.SugaredLogger) *Handler {
	loger.Debugln("Create new API handler")
	return &Handler{
		service: service,
		loger:   loger,
	}
}

// postMetric : Обработка запроса для создания метрики
func (h *Handler) postMetric(w http.ResponseWriter, r *http.Request) {
	h.loger.Debugln("PostMetric handler start")
	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "metric")
	metricValue := chi.URLParam(r, "value")

	if metricName == "" {
		w.Header().Set("Content-Type", TextPlain)
		http.Error(w, "incorrect metric name", http.StatusNotFound)
		return
	}

	if metricValue == "" {
		w.Header().Set("Content-Type", TextPlain)
		http.Error(w, "incorrect metric value", http.StatusBadRequest)
		return
	}

	switch metricType {
	case "gauge":

		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			w.Header().Set("Content-Type", TextPlain)
			http.Error(w, "incorrect metric value", http.StatusBadRequest)
			return
		}
		err = h.service.SetGauge(metricName, value)
		if err != nil {
			w.Header().Set("Content-Type", TextPlain)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", TextPlain)
		w.WriteHeader(http.StatusOK)

	case "counter":

		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			w.Header().Set("Content-Type", TextPlain)
			http.Error(w, "incorrect metric value", http.StatusBadRequest)
			return
		}
		err = h.service.SetCounter(metricName, value)
		if err != nil {
			w.Header().Set("Content-Type", TextPlain)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", TextPlain)
		w.WriteHeader(http.StatusOK)

	default:

		w.Header().Set("Content-Type", TextPlain)
		http.Error(w, "incorrect metric type", http.StatusBadRequest)

	}
	h.loger.Debugln("PostMetric handler end")
}

// getMetric : Обработка запроса для получения значения метрики
func (h *Handler) getMetric(w http.ResponseWriter, r *http.Request) {
	h.loger.Debugln("GetMetrics handler start")
	var answer string
	w.Header().Set("Content-Type", TextPlain)

	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "metric")

	switch metricType {
	case "gauge":

		value, exists := h.service.GetGauge(metricName)
		if !exists {
			http.Error(w, "unknown metric", http.StatusNotFound)
			return
		}
		answer = strconv.FormatFloat(value, 'f', -1, 64)

	case "counter":

		value, exists := h.service.GetCounter(metricName)
		if !exists {
			http.Error(w, "unknown metric", http.StatusNotFound)
			return
		}
		answer = fmt.Sprintf("%d", value)

	default:

		w.Header().Set("Content-Type", TextPlain)
		http.Error(w, "incorrect metric type", http.StatusBadRequest)

	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(answer))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.loger.Debugln("GetMetrics handler end")
}

// indexHandler : Обработка запроса для получения списка метрик на html странице
func (h *Handler) indexHandler(w http.ResponseWriter, r *http.Request) {
	h.loger.Debugln("Index handler start")
	var result string

	for metric, value := range h.service.GetCounters() {
		result += fmt.Sprintf("%s: %d\n", metric, value)
	}

	m := h.service.GetGauges()

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		result += fmt.Sprintf("%s: %.3f\n", k, m[k])
	}

	htmlString := "<html><head><title>Metrics</title></head><body><pre>" +
		html.EscapeString(result) +
		"</pre></body></html>"

	w.Header().Set("Content-Type", TextHTML)
	w.WriteHeader(http.StatusOK)

	if _, err = w.Write([]byte(htmlString)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.loger.Debugln("Index handler end")
}

// postJSONMetric : Обработка запроса на создание метрики в JSON формате
func (h *Handler) postJSONMetric(w http.ResponseWriter, r *http.Request) {
	h.loger.Debugln("PostJSONMetric handler start")

	var metric models.Metrics
	var buf bytes.Buffer
	reqID := middleware.GetReqID(r.Context())

	_, err = buf.ReadFrom(r.Body)
	if err != nil {
		h.jsonError(w, err.Error(), http.StatusBadRequest, reqID)
		return
	}

	h.loger.Debugln("Request body:", buf.String())

	if err = json.Unmarshal(buf.Bytes(), &metric); err != nil {
		h.jsonError(w, err.Error(), http.StatusBadRequest, reqID)
		return
	}

	switch metric.MType {
	case "gauge":

		err = h.service.SetGauge(metric.ID, *metric.Value)
		if err != nil {
			h.jsonError(w, "internal server error", http.StatusInternalServerError, reqID)
			return
		}

	case "counter":

		err = h.service.SetCounter(metric.ID, *metric.Delta)
		if err != nil {
			h.jsonError(w, "internal server error", http.StatusInternalServerError, reqID)
			return
		}

	default:

		h.jsonError(w, "incorrect metric type", http.StatusBadRequest, reqID)
		return
	}

	h.withJSON(w, metric, http.StatusOK, reqID)
	h.loger.Debugln("PostJSONMetric handler end")
}

// getJSONMetric : Обработка запроса для получения значения метрики в JSON формате
func (h *Handler) getJSONMetric(w http.ResponseWriter, r *http.Request) {
	h.loger.Debugln("GetJSONMetric handler start")

	var metric models.Metrics
	var buf bytes.Buffer
	reqID := middleware.GetReqID(r.Context())

	_, err = buf.ReadFrom(r.Body)
	if err != nil {
		h.jsonError(w, err.Error(), http.StatusBadRequest, reqID)
		return
	}

	h.loger.Debugln("Request body:", buf.String())

	if err = json.Unmarshal(buf.Bytes(), &metric); err != nil {
		h.jsonError(w, err.Error(), http.StatusBadRequest, reqID)
		return
	}

	switch metric.MType {
	case "gauge":

		value, exists := h.service.GetGauge(metric.ID)
		if !exists {
			h.loger.Debugln("GetJSONMetric gauge not found")
			h.jsonError(w, "unknown metric", http.StatusNotFound, reqID)
			return
		}
		metric.Value = &value

	case "counter":

		value, exists := h.service.GetCounter(metric.ID)
		if !exists {
			h.jsonError(w, "unknown metric", http.StatusNotFound, reqID)
			return
		}
		metric.Delta = &value

	default:

		h.jsonError(w, "incorrect metric type", http.StatusBadRequest, reqID)
		return

	}
	h.withJSON(w, metric, http.StatusOK, reqID)
	h.loger.Debugln("GetJSONMetric handler end")
}

// withJSON : Отправка ответа в JSON формате
func (h *Handler) withJSON(w http.ResponseWriter, v any, status int, reqID string) {
	h.loger.Debugln("withJSON util start")

	w.Header().Add("Content-Type", ApplicationJSON)
	w.WriteHeader(status)
	if err = json.NewEncoder(w).Encode(v); err != nil {
		h.jsonError(w, "failed to encode", http.StatusInternalServerError, reqID)
	}
	h.loger.Debugln("withJSON util end")
}

// jsonError : Обработка ошибок в JSON формате
func (h *Handler) jsonError(w http.ResponseWriter, error string, code int, reqID string) {
	h.loger.Debugln("JSON Error util start")

	var resp []byte
	w.WriteHeader(code)
	w.Header().Set("Content-Type", ApplicationJSON)
	answer := models.ErrorResponse{
		RequestID: reqID,
		Error:     error,
	}
	resp, err = json.Marshal(answer)
	if err != nil {
		h.loger.Errorf("Error marshalling response: %v", err)
		return
	}
	_, err = w.Write(resp)
	if err != nil {
		h.loger.Errorf("Error writing response: %v", err)
		return
	}
	h.loger.Debugln("JSON Error util end")
}
