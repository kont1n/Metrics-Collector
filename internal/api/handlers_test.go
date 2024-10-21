package api

import (
	"context"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"

	"Metrics-Collector/internal/service"
	"Metrics-Collector/internal/storage"
)

func TestEmptyMetricName(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/update/gauge/", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("type", "gauge")

	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	h := NewHandler(nil, sugar)
	h.postMetric(w, r)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestEmptyMetricValue(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/update/gauge/test", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("type", "gauge")
	rctx.URLParams.Add("metric", "test")

	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	h := NewHandler(nil, sugar)
	h.postMetric(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestIncorrectMetricType(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/update/test/test/1", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("type", "test")
	rctx.URLParams.Add("metric", "test")
	rctx.URLParams.Add("value", "1")

	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	h := NewHandler(nil, sugar)
	h.postMetric(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestIncorrectMetricValue(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/update/gauge/test/test", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("type", "gauge")
	rctx.URLParams.Add("metric", "test")
	rctx.URLParams.Add("value", "test")

	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	h := NewHandler(nil, sugar)
	h.postMetric(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestIncorrectMetricValue2(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/update/counter/test/test", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("type", "counter")
	rctx.URLParams.Add("metric", "test")
	rctx.URLParams.Add("value", "test")

	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	h := NewHandler(nil, sugar)
	h.postMetric(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPostGaugeMetric(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/update/gauge/test/1.0", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("type", "gauge")
	rctx.URLParams.Add("metric", "test")
	rctx.URLParams.Add("value", "1.0")

	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	store := storage.NewStore(sugar)
	srv := service.NewService(store, sugar)
	h := NewHandler(srv, sugar)
	h.postMetric(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
	value, _ := store.MemStorage.GetGauge("test")
	assert.Equal(t, 1.0, value)
}

func TestPostCounterMetric(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/update/counter/test/1", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("type", "counter")
	rctx.URLParams.Add("metric", "test")
	rctx.URLParams.Add("value", "1")

	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	store := storage.NewStore(sugar)
	srv := service.NewService(store, sugar)
	h := NewHandler(srv, sugar)
	h.postMetric(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
	value, _ := store.MemStorage.GetCounter("test")
	assert.Equal(t, int64(1), value)
}

func TestGetIncorrectMetricType(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/value/test/test", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("type", "test")
	rctx.URLParams.Add("metric", "test")

	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	h := NewHandler(nil, sugar)
	h.getMetric(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUnknownGaugeMetric(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/value/gauge/test", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("type", "gauge")
	rctx.URLParams.Add("metric", "test")

	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	store := storage.NewStore(sugar)
	srv := service.NewService(store, sugar)
	h := NewHandler(srv, sugar)
	h.getMetric(w, r)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUnknownCounterMetric(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/value/counter/test", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("type", "counter")
	rctx.URLParams.Add("metric", "test")

	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	store := storage.NewStore(sugar)
	srv := service.NewService(store, sugar)
	h := NewHandler(srv, sugar)
	h.getMetric(w, r)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetGaugeMetrics(t *testing.T) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	store := storage.NewStore(sugar)
	srv := service.NewService(store, sugar)
	h := NewHandler(srv, sugar)

	store.MemStorage.SetGauge("test", 1.0)

	r := httptest.NewRequest(http.MethodGet, "/value/gauge/test", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("type", "gauge")
	rctx.URLParams.Add("metric", "test")

	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	h.getMetric(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "1", w.Body.String())
}

func TestGetCounterMetrics(t *testing.T) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	store := storage.NewStore(sugar)
	srv := service.NewService(store, sugar)
	h := NewHandler(srv, sugar)

	store.MemStorage.SetCounter("test", int64(1))

	r := httptest.NewRequest(http.MethodGet, "/value/counter/test", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("type", "counter")
	rctx.URLParams.Add("metric", "test")

	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	h.getMetric(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "1", w.Body.String())
}
