package service

import (
	"go.uber.org/zap"

	"Metrics-Collector/internal/storage"
)

type Service struct {
	store *storage.Store
	loger *zap.SugaredLogger
}

func NewService(store *storage.Store, loger *zap.SugaredLogger) *Service {
	return &Service{
		store: store,
		loger: loger,
	}
}

func (s Service) SetGauge(name string, value float64) {
	s.loger.Debugln("SetGauge")
	s.store.MemStorage.SetGauge(name, value)
}

func (s Service) SetCounter(name string, value int64) {
	s.loger.Debugln("SetCounter")
	s.store.MemStorage.SetCounter(name, value)
}

func (s Service) GetGauge(name string) (value float64, exists bool) {
	s.loger.Debugln("GetGauge")
	return s.store.MemStorage.GetGauge(name)
}

func (s Service) GetCounter(name string) (value int64, exists bool) {
	s.loger.Debugln("GetCounter")
	return s.store.MemStorage.GetCounter(name)
}

func (s Service) GetCounters() map[string]int64 {
	s.loger.Debugln("GetCounters")
	return s.store.MemStorage.GetCounters()
}

func (s Service) GetGauges() map[string]float64 {
	s.loger.Debugln("GetGauges")
	return s.store.MemStorage.GetGauges()
}
