package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/keruch/tfs-go-hw/trading_robot/data_receiver/internal/domain"
	"github.com/keruch/tfs-go-hw/trading_robot/data_receiver/internal/service"
	"github.com/keruch/tfs-go-hw/trading_robot/data_receiver/pkg/utils"
)

type Handler interface {
	SendTickers(ctx context.Context, tickers []string, u *url.URL) error
}

type ReceiverHandler struct {
	service service.Service
	client  *utils.RetryableClient
}

func NewProcessorHandler(service service.Service) Handler {
	return &ReceiverHandler{
		service: service,
		client: utils.NewRetryableClient(&http.Client{
			Timeout: 1 * time.Second,
		}),
	}
}

func (rh *ReceiverHandler) SendTickers(ctx context.Context, tickers []string, u *url.URL) error {
	tickersChan := rh.service.Subscribe(ctx, tickers)

	go func() {
		for ticker := range tickersChan {
			// TODO: beautify timestamps
			//fmt.Println(time.UnixMilli(int64(ticker.Time)))

			err := rh.sendData(ticker, u)
			if err != nil {
				panic(err)
			}
		}
	}()

	return nil
}

func (rh *ReceiverHandler) sendData(ticker domain.Ticker, u *url.URL) error {
	tickerData, err := json.Marshal(ticker)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(tickerData)
	req, err := http.NewRequest(http.MethodPost, u.String(), reader)
	if err != nil {
		return err
	}

	err = rh.client.SendRequestWithoutTimeout(req)
	if err != nil {
		return err
	}

	return nil
}
