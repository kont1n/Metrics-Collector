package storage

import "sync"

type MemStorage struct {
	gaugeMetrics   map[string]float64
	counterMetrics map[string]int64
	mu             sync.RWMutex
}

type Storage interface {
	SetGauge(key string, value float64)
	SetCounter(key string, value int64)
	GetGauge(key string) float64
	GetCounter(key string) int64
	GetGauges() map[string]float64
	GetCounters() map[string]int64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gaugeMetrics:   make(map[string]float64),
		counterMetrics: make(map[string]int64),
	}
}

func (m *MemStorage) SetGauge(key string, value float64) {
	m.mu.Lock()
	m.gaugeMetrics[key] = value
	m.mu.Unlock()
}

func (m *MemStorage) SetCounter(key string, value int64) {
	m.mu.Lock()
	m.counterMetrics[key] += value
	m.mu.Unlock()
}

func (m *MemStorage) GetGauge(key string) (float64, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	value, exists := m.gaugeMetrics[key]
	if exists {
		return value, true
	} else {
		return value, false
	}
}

func (m *MemStorage) GetCounter(key string) (int64, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	value, exists := m.counterMetrics[key]
	if exists {
		return value, true
	} else {
		return value, false
	}
}

func (m *MemStorage) GetGauges() map[string]float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.gaugeMetrics
}

func (m *MemStorage) GetCounters() map[string]int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.counterMetrics
}
