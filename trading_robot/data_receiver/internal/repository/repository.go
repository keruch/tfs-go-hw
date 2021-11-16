package repository

import (
	"context"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/keruch/tfs-go-hw/trading_robot/data_receiver/pkg/utils"
)

type Repository interface {
	GetTickersData(ctx context.Context, tickers []string) (<-chan []byte, error)
	Shutdown() error
}

const (
	scheme = "wss"
	path   = "/ws/v1"
	host   = "demo-futures.kraken.com"

	maxRetries = 13

	eventType = "subscribe"
	feedType  = "ticker"
)

type KrakenExchange struct {
	conn *websocket.Conn
}

func NewKrakenExchange() (Repository, error) {
	u := url.URL{
		Scheme: scheme,
		Host:   host,
		Path:   path,
	}

	conn, _, err := utils.CreateRetryableDial(u, maxRetries)
	if err != nil {
		return nil, err
	}

	return &KrakenExchange{
		conn: conn,
	}, nil
}

func (ke *KrakenExchange) GetTickersData(ctx context.Context, tickers []string) (<-chan []byte, error) {
	err := ke.requestData(tickers)
	if err != nil {
		return nil, err
	}

	out := make(chan []byte)

	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				_, rawTicker, err := ke.conn.ReadMessage()
				if err != nil {
					panic(err)
				}

				out <- rawTicker
			}
		}
	}()

	return out, nil
}

func (ke *KrakenExchange) Shutdown() error {
	return ke.conn.Close()
}

type request struct {
	Event      string   `json:"event"`
	Feed       string   `json:"feed"`
	ProductIDs []string `json:"product_ids"`
}

func newRequest() request {
	return request{
		Event: eventType,
		Feed:  feedType,
	}
}

func (r *request) addProducts(ticker []string) {
	r.ProductIDs = append(r.ProductIDs, ticker...)
}

func (ke *KrakenExchange) requestData(tickers []string) error {
	r := newRequest()
	r.addProducts(tickers)

	err := ke.conn.WriteJSON(r)
	if err != nil {
		return err
	}
	_, _, err = ke.conn.ReadMessage()
	if err != nil {
		return err
	}

	_, _, err = ke.conn.ReadMessage()
	if err != nil {
		return err
	}

	return nil
}
