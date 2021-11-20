package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/keruch/tfs-go-hw/trading_robot/internal/domain"
	"github.com/keruch/tfs-go-hw/trading_robot/internal/exchange"
	"github.com/keruch/tfs-go-hw/trading_robot/pkg/log"
)

func init() {
	os.Setenv("PUBLIC_KEY", "OLlMf8FawHE0FC4L7wZmvpxDEZcz184C0DbG6WgX3+ZoLb4s9EmtI9VR")
	os.Setenv("PRIVATE_KEY", "tm+El0L45IDVJLoy05xRWf8Bu7IcyRpsrXwYe5GPtuLlKzcxE3rluKNf6oggf5E0jEUmgtGUo8WXywshBSEToZxM")
}

func main() {
	// PUBLIC_KEY and PRIVATE_KEY env variables have to be set
	logger := log.NewLogger()
	ex, err := exchange.NewKrakenExchange(logger)
	if err != nil {
		logger.Errorf("NewKrakenService: %s", err)
		return
	}

	out := ex.Listen()

	// POST /pair - to change pair
	// POST /index/{params} - for index
	// POST /operations/shutdown - for shutdown
	// (?) POST /operations/start - for start
	// (?) POST /operations/stop - for stop

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)

	router := mux.NewRouter()
	router.Methods(http.MethodPost).PathPrefix(domain.SubscribePair).HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			// TODO: add error if user will call subscribe twice
			vars := mux.Vars(request)
			ticker := vars[domain.PairVar]
			err := ex.SubscribePairs(ticker)
			if err != nil {
				logger.Errorf("%s endpoint: %s", domain.SubscribePair, err)
				writer.WriteHeader(http.StatusBadRequest)
			}
		})

	router.Methods(http.MethodPost).PathPrefix(domain.UnsubscribePair).HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			// TODO: add error if user will call subscribe twice
			vars := mux.Vars(request)
			ticker := vars[domain.PairVar]
			err := ex.UnsubscribePairs(ticker)
			if err != nil {
				logger.Errorf("%s endpoint: %s", domain.UnsubscribePair, err)
				writer.WriteHeader(http.StatusBadRequest)
			}
		})
	//router.Methods(http.MethodPost).PathPrefix(kraken.IndexEndpoint).HandlerFunc(h.GetOrders)
	router.Methods(http.MethodPost).PathPrefix(domain.OperationsEndpoint + domain.ShutdownOperation).HandlerFunc(
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
		if err := srv.ListenAndServe(); err != nil {
			logger.Infof("ListenAndServe: %s", err)
		}
	}()

	go func() {
		for data := range out {
			logger.Trace(string(data))
		}
	}()

	<-shutdown

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), domain.WaitTime)
	defer shutdownCancel()

	go func() {
		srvCtx, srvCancel := context.WithCancel(shutdownCtx)
		if err := srv.Shutdown(srvCtx); err != nil {
			logger.Errorf("Server shutdown: %s", err)
		}
		srvCancel()
		logger.Infof("Server closed")
	}()

	go func() {
		if err := ex.Shutdown(); err != nil {
			logger.Errorf("Kraken exchange shutdown: %s", err)
		}
		logger.Infof("Kraken exchange closed")
	}()

	<-shutdownCtx.Done()

	logger.Infof("Trading robot closed")
}
