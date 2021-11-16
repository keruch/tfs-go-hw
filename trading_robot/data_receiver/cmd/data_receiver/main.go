package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/signal"

	"github.com/keruch/tfs-go-hw/trading_robot/data_receiver/internal/handlers"
	"github.com/keruch/tfs-go-hw/trading_robot/data_receiver/internal/repository"
	"github.com/keruch/tfs-go-hw/trading_robot/data_receiver/internal/service"
)

const (
	serverScheme = "http"
	serverHost   = "localhost"
	serverPort   = ":8081"
)

func main() {
	var tickers = []string{"PI_ETHUSD"}

	r, err := repository.NewKrakenExchange()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer r.Shutdown()

	ctx, cancel := context.WithCancel(context.Background())

	s := service.NewReceiverService(r)

	h := handlers.NewProcessorHandler(s)

	u := &url.URL{
		Scheme: serverScheme,
		Path:   serverHost + serverPort,
	}

	err = h.SendTickers(ctx, tickers, u)
	if err != nil {
		fmt.Println(err)
		return
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c
	cancel()
}
