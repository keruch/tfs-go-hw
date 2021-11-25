package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/keruch/tfs-go-hw/trading_robot/config"
	"github.com/keruch/tfs-go-hw/trading_robot/internal/candles"
	"github.com/keruch/tfs-go-hw/trading_robot/internal/domain"
	"github.com/keruch/tfs-go-hw/trading_robot/internal/exchange"
	"github.com/keruch/tfs-go-hw/trading_robot/internal/indicator"
	"github.com/keruch/tfs-go-hw/trading_robot/internal/processor"
	"github.com/keruch/tfs-go-hw/trading_robot/internal/repository"
	"github.com/keruch/tfs-go-hw/trading_robot/pkg/log"
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
		logger.Fatal(err)
	}

	// setup strategy
	strat := SetupStrategy()
	logger.Info("Setup strategy")

	// setup exchange
	ex, err := exchange.NewKrakenExchange(logger)
	if err != nil {
		logger.Fatal(err)
	}
	defer ex.CloseConnection()
	logger.Info("Setup exchange")

	// setup repository
	repo, err := repository.NewPostgreSQLPool(config.GetDatabaseURL(), logger)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("Setup repository")

	// setup orders processor
	proc := processor.NewOrdersProcessor(strat, repo, ex, logger)
	logger.Info("Setup processor")

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	// get prices from exchange
	prices := ex.GetPrices(ctx)

	// generate candles from prices
	wg.Add(1)
	candle := candles.GenerateCandles(prices, config.GetPeriod(), &wg)

	//skip the first candle (she is invalid)
	//<-candle

	// run processor
	wg.Add(1)
	proc.ProcessCandles(candle, &wg)
	logger.Info("Processing candles...")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// POST /pair - to change pair
	// POST /indicator/{params} - for indicator
	// POST /operations/shutdown - for shutdown
	// (?) POST /operations/start - for start
	// (?) POST /operations/stop - for stop

	router := mux.NewRouter()
	router.Methods(http.MethodPost).PathPrefix(domain.SubscribePair).HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			vars := mux.Vars(request)
			ticker := vars[domain.PairVar]
			err = ex.SubscribePairs(ticker)
			if err != nil {
				logger.Errorf("%s endpoint: %s", domain.SubscribePair, err)
				writer.WriteHeader(http.StatusBadRequest)
			}
		})

	router.Methods(http.MethodPost).PathPrefix(domain.UnsubscribePair).HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			vars := mux.Vars(request)
			ticker := vars[domain.PairVar]
			err = ex.UnsubscribePairs(ticker)
			if err != nil {
				logger.Errorf("%s endpoint: %s", domain.UnsubscribePair, err)
				writer.WriteHeader(http.StatusBadRequest)
			}
		})

	router.Methods(http.MethodPost).PathPrefix(domain.ShutdownOperation).HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			shutdown <- os.Interrupt
		})

	srv := &http.Server{
		Handler:      router,
		Addr:         ":8091",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		logger.Info("Setup server")
		if err = srv.ListenAndServe(); err != nil {
			logger.Infof("ListenAndServe: %s", err)
		}
	}()

	<-shutdown
	ctxTimeout, cancelTimeout := context.WithTimeout(context.Background(), 3*time.Second)

	go func() {
		if err = srv.Shutdown(ctxTimeout); err != nil {
			logger.Fatal(err)
		}
		cancelTimeout()
		logger.Info("Server done")
	}()

	cancel()
	wg.Wait()

	<-ctxTimeout.Done()

	logger.Infof("Trading robot close.")
}
