package collector

import (
	"fmt"
	"math/rand/v2"
	"net/http"
	"runtime"
	"sync"
	"time"
)

type Agent struct {
	ServerURL      string
	PollInterval   time.Duration
	ReportInterval time.Duration
	PollCount      int64
	Metrics        map[string]float64
	wg             sync.WaitGroup
	mu             sync.Mutex
}

func NewAgent(serverURL string, pollInterval time.Duration, reportInterval time.Duration) *Agent {
	return &Agent{
		ServerURL:      serverURL,
		PollInterval:   pollInterval,
		ReportInterval: reportInterval,
		PollCount:      0,
		Metrics:        make(map[string]float64),
	}
}

func (a *Agent) Run() {
	a.wg.Add(2)
	go a.Poll()
	go a.Report()
	a.wg.Wait()
}

func (a *Agent) Poll() {
	defer a.wg.Done()
	for {
		time.Sleep(a.PollInterval)

		a.mu.Lock()
		a.PollCount++
		runtimeMetrics := collectedGauges()
		runtimeMetrics["RandomValue"] = rand.Float64()
		a.Metrics = runtimeMetrics
		a.mu.Unlock()
	}
}

func (a *Agent) Report() {
	defer a.wg.Done()
	for {
		time.Sleep(a.ReportInterval)

		a.mu.Lock()
		url := fmt.Sprintf("%s/%s/%s/%d", a.ServerURL, "counter", "PollCount", a.PollCount)
		sendMetrics(url)
		for metric, value := range a.Metrics {
			url = fmt.Sprintf("%s/%s/%s/%f", a.ServerURL, "gauge", metric, value)
			sendMetrics(url)
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

func sendMetrics(url string) {
	fmt.Println("Send metric:", url)
	response, err := http.Post(url, "text/plain", nil)
	if err != nil {
		fmt.Println("Error sending metrics:", err.Error())
		return
	}
	fmt.Println("Status code:", response.StatusCode)
	err = response.Body.Close()
	if err != nil {
		fmt.Println("Error closing response body:", err.Error())
		return
	}
}
