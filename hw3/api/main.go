package main

import (
	"context"
	"os"
	"os/signal"
	"tfs-go-hw/hw3/internal/domain/generator"
	"tfs-go-hw/hw3/internal/repository"
	"tfs-go-hw/hw3/internal/services"
	"time"

	"tfs-go-hw/hw3/pkg/log"
)

var tickers = []string{"AAPL", "SBER", "NVDA", "TSLA"}

func main() {
	logger := log.NewLogger()
	logger.Infof("Starting service.")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		oscall := <-c
		logger.Infof("Handled %v signal. Shutting down.", oscall)
		cancel()
	}()

	pg := generator.NewPricesGenerator(generator.Config{
		Factor:  10,
		Delay:   time.Millisecond * 10,
		Tickers: tickers,
	})

	repo := repository.NewGeneratorData(pg)

	cs := services.NewCandlesService(repo, logger)
	err := cs.GenerateCandles(ctx)
	if err != nil {
		logger.Errorf("Error while generating candles: %v", err)
		return
	}
}
