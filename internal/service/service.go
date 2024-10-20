package service

import (
	"go.uber.org/zap"

	"Metrics-Collector/internal/storage"
)

type Service struct {
	store *storage.Store
	loger *zap.SugaredLogger
	Services
}

type Services interface {
	SetGauge(name string, value float64) error
	SetCounter(name string, value int64) error
	GetGauge(name string) (value float64, exists bool)
	GetCounter(name string) (value int64, exists bool)
	GetCounters() map[string]int64
	GetGauges() map[string]float64
}

func NewService(store *storage.Store, loger *zap.SugaredLogger) *Service {
	loger.Debugln("Create new service")
	return &Service{
		store: store,
		loger: loger,
	}
}

func (s Service) SetGauge(name string, value float64) error {
	s.loger.Debugln("SetGauge")
	s.store.MemStorage.SetGauge(name, value)
	return nil
}

func (s Service) SetCounter(name string, value int64) error {
	s.loger.Debugln("SetCounter")
	s.store.MemStorage.SetCounter(name, value)
	return nil
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
