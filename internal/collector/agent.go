package collector

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"runtime"
	"strconv"
	"sync"
	"time"

	"Metrics-Collector/internal/models"
)

const (
	TextPlain       = "text/plain"
	ApplicationJSON = "application/json"
)

type Agent struct {
	ServerURL      string
	PollInterval   time.Duration
	ReportInterval time.Duration
	PollCount      int64
	Metrics        map[string]float64
	wg             sync.WaitGroup
	mu             sync.Mutex
	log            *slog.Logger
}

func NewAgent(serverURL string, pollInterval time.Duration, reportInterval time.Duration, log *slog.Logger) *Agent {
	return &Agent{
		ServerURL:      serverURL,
		PollInterval:   pollInterval,
		ReportInterval: reportInterval,
		PollCount:      0,
		Metrics:        make(map[string]float64),
		log:            log,
	}
}

func (a *Agent) Run() {
	a.log.Debug("Agent run")
	a.wg.Add(2)
	go a.Poll()
	//go a.Report()
	go a.ReportJSON()
	a.wg.Wait()
}

func (a *Agent) Poll() {
	a.log.Debug("Poll start")
	defer a.wg.Done()
	for {
		time.Sleep(a.PollInterval)

		a.log.Debug("Collect start")
		a.mu.Lock()
		a.PollCount++
		runtimeMetrics := collectedGauges()
		runtimeMetrics["RandomValue"] = rand.Float64()
		a.Metrics = runtimeMetrics
		a.log.Debug("Metrics collected",
			slog.Int64("PollCount", a.PollCount),
			slog.Any("metrics", a.Metrics))
		a.mu.Unlock()
	}
}

func (a *Agent) Report() {
	a.log.Debug("Report start")
	defer a.wg.Done()
	for {
		time.Sleep(a.ReportInterval)

		a.mu.Lock()
		url := fmt.Sprintf("%s/%s/%s/%d", a.ServerURL, "counter", "PollCount", a.PollCount)
		sendMetrics(url, a.log)
		for metric, value := range a.Metrics {
			url = fmt.Sprintf("%s/%s/%s/%f", a.ServerURL, "gauge", metric, value)
			sendMetrics(url, a.log)
		}
		a.PollCount = 0
		a.mu.Unlock()
	}
}

func (a *Agent) ReportJSON() {
	a.log.Debug("ReportJSON start")
	var counterData, gaugeData models.Metrics
	defer a.wg.Done()
	for {
		time.Sleep(a.ReportInterval)

		a.mu.Lock()
		counterData = models.Metrics{
			ID:    "PollCount",
			MType: "counter",
			Delta: &a.PollCount,
		}
		sendJSONMetrics(a.ServerURL, counterData, a.log)
		for metric, value := range a.Metrics {
			gaugeData = models.Metrics{
				ID:    metric,
				MType: "gauge",
				Value: &value,
			}
			sendJSONMetrics(a.ServerURL, gaugeData, a.log)
		}
		a.PollCount = 0
		a.mu.Unlock()
	}
}

func collectedGauges() map[string]float64 {
	var m runtime.MemStats
	metrics := make(map[string]float64)

	metrics["Alloc"] = float64(m.Alloc)
	metrics["BuckHashSys"] = float64(m.BuckHashSys)
	metrics["Frees"] = float64(m.Frees)
	metrics["GCCPUFraction"] = m.GCCPUFraction
	metrics["GCSys"] = float64(m.GCSys)
	metrics["HeapAlloc"] = float64(m.HeapAlloc)
	metrics["HeapIdle"] = float64(m.HeapIdle)
	metrics["HeapInuse"] = float64(m.HeapInuse)
	metrics["HeapObjects"] = float64(m.HeapObjects)
	metrics["HeapReleased"] = float64(m.HeapReleased)
	metrics["HeapSys"] = float64(m.HeapSys)
	metrics["LastGC"] = float64(m.LastGC)
	metrics["Lookups"] = float64(m.Lookups)
	metrics["MCacheInuse"] = float64(m.MCacheInuse)
	metrics["MCacheSys"] = float64(m.MCacheSys)
	metrics["MSpanInuse"] = float64(m.MSpanInuse)
	metrics["MSpanSys"] = float64(m.MSpanSys)
	metrics["Mallocs"] = float64(m.Mallocs)
	metrics["NextGC"] = float64(m.NextGC)
	metrics["OtherSys"] = float64(m.OtherSys)
	metrics["PauseTotalNs"] = float64(m.PauseTotalNs)
	metrics["StackInuse"] = float64(m.StackInuse)
	metrics["StackSys"] = float64(m.StackSys)
	metrics["Sys"] = float64(m.Sys)
	metrics["TotalAlloc"] = float64(m.TotalAlloc)

	return metrics
}

func sendMetrics(url string, log *slog.Logger) {
	log.Debug("Send metric",
		slog.String("url", url))

	response, err := http.Post(url, TextPlain, nil)
	if err != nil {
		log.Error("Error sending metrics:",
			slog.String("error", err.Error()))
		return
	}
	log.Debug("Response received",
		slog.String("Status code", strconv.Itoa(response.StatusCode)))
	err = response.Body.Close()
	if err != nil {
		log.Error("Error closing response body:",
			slog.String("error", err.Error()))
		return
	}
	log.Debug("Send metric end")
}

func sendJSONMetrics(url string, metric models.Metrics, log *slog.Logger) {
	log.Debug("Send JSON metric",
		slog.String("url", url))

	body, err := json.Marshal(metric)
	if err != nil {
		log.Error("Error marshalling song",
			slog.String("error", err.Error()))
		return
	}
	log.Debug("Metric info",
		slog.String("metric", string(body)))

	buf := bytes.NewBuffer(body)

	response, err := http.Post(url, ApplicationJSON, buf)
	if err != nil {
		log.Error("Error sending metrics",
			slog.String("error", err.Error()))
		return
	}
	log.Debug("Response received",
		slog.String("Status code", strconv.Itoa(response.StatusCode)))
	err = response.Body.Close()
	if err != nil {
		log.Error("Error closing response body",
			slog.String("error", err.Error()))
		return
	}
	log.Debug("Send JSON metric end")
}
