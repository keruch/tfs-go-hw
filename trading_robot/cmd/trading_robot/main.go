package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/keruch/tfs-go-hw/trading_robot/config"
	"github.com/keruch/tfs-go-hw/trading_robot/internal/exchange"
	"github.com/keruch/tfs-go-hw/trading_robot/internal/indicator"
	"github.com/keruch/tfs-go-hw/trading_robot/internal/processor"
	"github.com/keruch/tfs-go-hw/trading_robot/internal/repository"
	"github.com/keruch/tfs-go-hw/trading_robot/internal/router"
	"github.com/keruch/tfs-go-hw/trading_robot/pkg/log"
	"github.com/keruch/tfs-go-hw/trading_robot/pkg/tg"
)

func SetupStrategy() indicator.Strategy {
	alphaFunc := func(p int) float64 {
		return 2 / float64(p+1)
	}
	period := 100
	//macd := NewMACDEvaluator(12, 26, 9, alphaFunc)
	ema := indicator.NewEMAEvaluator(period, alphaFunc)
	//NewMACDStrategy(macd)
	return indicator.NewStrategiesComposition(indicator.NewEMAStrategy(ema))
}

func main() {
	// setup logger and config
	logger := log.NewLogger()
	err := config.SetupConfig()
	if err != nil {
		logger.Fatalf("Setup config failed: %s", err)
	}
	logger.Info("Setup config")

	// setup strategy
	strategy := SetupStrategy()
	logger.Info("Setup strategy")

	// setup exchange
	kraken, err := exchange.NewKrakenExchange(logger)
	if err != nil {
		logger.Fatalf("Setup exchange failed: %s", err)
	}
	defer kraken.CloseConnection()
	logger.Info("Setup exchange")

	// setup repository
	repo, err := repository.NewPostgreSQLPool(config.GetDatabaseURL(), logger)
	if err != nil {
		logger.Fatalf("Setup repository failed: %s", err)
	}
	logger.Info("Setup repository")

	// setup telegram bot
	telegram, err := tg.NewTelegramBot(config.GetTelegramBotToken(), logger)
	if err != nil {
		logger.Fatalf("Setup telegram failed: %s", err)
	}
	logger.Info("Setup telegram bot")

	// setup orders processor
	proc := processor.NewOrdersProcessor(strategy, repo, kraken, telegram, logger)
	logger.Info("Setup processor")

	// setup router
	r := router.NewRouter(kraken, proc, logger)
	logger.Info("Setup router")

	// setup server
	srv := &http.Server{
		Handler:      r,
		Addr:         config.GetServerAddress(),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	logger.Info("Setup server")

	botCtx, botShutdown := context.WithCancel(context.Background())

	// start tg bot
	//go telegram.Serve(botCtx)

	// start processing
	var shutdownWait sync.WaitGroup
	proc.StartTradingBotProcessor(botCtx, &shutdownWait)

	// setup signals handler
	shutdownSig := make(chan os.Signal, 1)
	signal.Notify(shutdownSig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-shutdownSig

		// give 5 seconds to shutdown
		forceShutdownCtx, forceShutdown := context.WithTimeout(botCtx, time.Second*5)
		go func() {
			<-forceShutdownCtx.Done()
			if forceShutdownCtx.Err() == context.DeadlineExceeded {
				logger.Fatal("graceful shutdown timed out, forcing exit")
			}
			forceShutdown()
		}()

		// Trigger graceful shutdown
		err = srv.Shutdown(forceShutdownCtx)
		if err != nil {
			logger.Fatal(err)
		}
		logger.Info("Server done")

		botShutdown()
		err = kraken.CloseConnection()
		if err != nil {
			logger.Fatal(err)
		}
		logger.Info("WS exchange connection done")
	}()

	if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("ListenAndServe: %s", err)
	}
	
	shutdownWait.Wait()

	logger.Infof("Trading robot close")
}
