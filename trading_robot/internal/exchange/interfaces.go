package exchange

import (
	"github.com/keruch/tfs-go-hw/trading_robot/internal/domain"
	"github.com/keruch/tfs-go-hw/trading_robot/internal/exchange/kraken"
)

type TickersGetter interface {
	GetTickers() <-chan domain.Ticker
}

type Exchange interface {
	Listen() <-chan []byte
	// Add get tickers
	SubscribePairs(pairs ...string) error
	UnsubscribePairs(pairs ...string) error
	Shutdown() error

	// TODO: cahnge orders to domain

	CreateOrder(order *kraken.SendOrder) (*kraken.ReceiveOrder, error)
	GetOrders(order *kraken.SendOrder) (*kraken.ReceiveOrder, error)
	DeleteOrder(order *kraken.SendOrder) (*kraken.ReceiveOrder, error)
}
