package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
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
		logger.Fatalf("Setup config failed: %s",err)
	}
	logger.Info("Setup config")

	// setup strategy
	strat := SetupStrategy()
	logger.Info("Setup strategy")

	// setup exchange
	ex, err := exchange.NewKrakenExchange(logger)
	if err != nil {
		logger.Fatalf("Setup exchange failed: %s",err)
	}
	defer ex.CloseConnection()
	logger.Info("Setup exchange")

	// setup repository
	repo, err := repository.NewPostgreSQLPool(config.GetDatabaseURL(), logger)
	if err != nil {
		logger.Fatalf("Setup repository failed: %s", err)
	}
	logger.Info("Setup repository")

	// setup telegram bot
	tgBot, err := tg.NewTelegramBot(config.GetTelegramBotToken(), logger)
	if err != nil {
		logger.Fatalf("Setup telegram failed: %s",err)
	}
	logger.Info("Setup telegram bot")

	// setup orders processor
	proc := processor.NewOrdersProcessor(strat, repo, ex, tgBot, logger)
	logger.Info("Setup processor")

	// setup signals handler
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// POST /pair - to change pair
	// POST /indicator/{params} - for indicator
	// POST /operations/shutdown - for shutdown
	// (?) POST /operations/start - for start
	// (?) POST /operations/stop - for stop

	// setup router
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

	router.Methods(http.MethodPost).PathPrefix(domain.SetQuantity).HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			vars := mux.Vars(request)
			quantity := vars[domain.ValueVal]
			quantityInt, err := strconv.Atoi(quantity)
			if err != nil {
				logger.Errorf("%s endpoint: %s", domain.SetQuantity, err)
				writer.WriteHeader(http.StatusBadRequest)
			}
			proc.SetTradingQuantity(quantityInt)
		})

	router.Methods(http.MethodPost).PathPrefix(domain.SetMultiplier).HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			vars := mux.Vars(request)
			multiplier := vars[domain.ValueVal]
			multiplierFloat, err := strconv.ParseFloat(multiplier, 32);
			if err != nil {
				logger.Errorf("%s endpoint: %s", domain.SetMultiplier, err)
				writer.WriteHeader(http.StatusBadRequest)
			}
			proc.SetPriceMultiplier(multiplierFloat)
		})

	router.Methods(http.MethodPost).PathPrefix(domain.ShutdownOperation).HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			shutdown <- syscall.SIGINT
		})

	// setup server
	srv := &http.Server{
		Handler:      router,
		Addr:         ":8091",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// start server
	go func() {
		logger.Info("Setup server")
		if err = srv.ListenAndServe(); err != nil {
			logger.Infof("ListenAndServe: %s", err)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())

	// start tg bot
	tgBot.Serve(ctx)

	var wg sync.WaitGroup
	// get prices from exchange
	prices := ex.GetPrices(ctx)

	// generate candles from prices
	wg.Add(1)
	candle := candles.GenerateCandles(prices, config.GetPeriod(), &wg)

	//skip the first candle (she is invalid)
	<-candle

	// run processor
	wg.Add(1)
	proc.ProcessCandles(candle, &wg)
	logger.Info("Processing candles...")

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
