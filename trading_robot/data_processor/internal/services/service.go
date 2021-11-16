package services

import (
	"github.com/keruch/tfs-go-hw/trading_robot/data_processor/internal/handlers"
	"github.com/keruch/tfs-go-hw/trading_robot/data_processor/internal/repository"
)

type Service interface {
	Start()
}

type DataProcessor struct {
	re
}

func NewService(repository *repository.Repository, handler *handlers.Handler) Service {
	return nil
}
