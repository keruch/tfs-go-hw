package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/keruch/tfs-go-hw/trading_robot/kraken/internal/domain"
	"github.com/keruch/tfs-go-hw/trading_robot/kraken/internal/handlers"
	"github.com/keruch/tfs-go-hw/trading_robot/kraken/internal/repository"
	"github.com/keruch/tfs-go-hw/trading_robot/kraken/internal/service"
)

func main() {
	r := repository.NewKrakenExchange()
	s := service.NewKrakenService(r, domain.PUBLIC_KEY)
	h := handlers.NewKrakenHandler(s)

	router := mux.NewRouter()
	router.Methods(http.MethodPost).PathPrefix(domain.OrderEndpoint).HandlerFunc(h.PostOrder)
	router.Methods(http.MethodGet).PathPrefix(domain.OrderEndpoint).HandlerFunc(h.GetOrders)
	router.Methods(http.MethodPut).PathPrefix(domain.OrderEndpoint).HandlerFunc(h.PutOrder)
	router.Methods(http.MethodDelete).PathPrefix(domain.OrderEndpoint).HandlerFunc(h.DeleteOrder)

	srv := &http.Server{
		Handler:      router,
		Addr:         ":8091",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
