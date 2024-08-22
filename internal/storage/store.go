package storage

type MemStorage struct {
	gaugeMetrics   map[string]float64
	counterMetrics map[string]int64
}

type Storage interface {
	SetGauge(key string, value float64)
	SetCounter(key string, value int64)
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gaugeMetrics:   make(map[string]float64),
		counterMetrics: make(map[string]int64),
	}
}

func (m *MemStorage) SetGauge(key string, value float64) {
	m.gaugeMetrics[key] = value
}

func (m *MemStorage) SetCounter(key string, value int64) {
	m.counterMetrics[key] += value
}
