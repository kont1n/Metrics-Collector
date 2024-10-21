package storage

import (
	"go.uber.org/zap"
)

type Store struct {
	MemStorage *MemStorage
	loger      *zap.SugaredLogger
}

func NewStore(loger *zap.SugaredLogger) *Store {
	loger.Debugln("Create new store")
	return &Store{
		MemStorage: NewMemStorage(),
		loger:      loger,
	}
}
