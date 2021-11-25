package kraken

import (
	"errors"
)

const (
	MaxRetries = 13

	Scheme = "https"
	Host   = "demo-futures.kraken.com"
	Path   = "/derivatives"

	WsScheme = "wss"
	WsPath   = "/ws/v1"
	WsQuery  = "chart"

	SubscribeEvent   Event = "subscribe"
	UnsubscribeEvent Event = "unsubscribe"

	Candles1mFeed Feed = "candles_trade_1m"
	TickerFeed    Feed = "ticker"
	TradesFeed    Feed = "trade"

	FeedType = TradesFeed

	CreateOrder OperationEndpoint = "/api/v3/sendorder"
	OpenOrders  OperationEndpoint = "/api/v3/openorders"
	EditOrder   OperationEndpoint = "/api/v3/editorder" // does not supported because of orderId filed
	CancelOrder OperationEndpoint = "/api/v3/cancelorder"

	Authent RequestHeader = "Authent"
	APIKey  RequestHeader = "APIKey"

	OrderID       QueryParam = "order_id"
	OrderType     QueryParam = "orderType"
	Symbol        QueryParam = "symbol"
	Side          QueryParam = "side"
	Size          QueryParam = "size"
	LimitPrice    QueryParam = "limitPrice"
	StopPrice     QueryParam = "stopPrice"
	TriggerSignal QueryParam = "triggerSignal"
	CliOrdId      QueryParam = "cliOrdId"
	ReduceOnly    QueryParam = "reduceOnly"
)

type (
	Event string
	Feed  string

	OperationEndpoint string

	RequestHeader = string

	QueryParam  = string
	QueryParams map[QueryParam]string

	OrderEvents struct {
		Type string `json:"type"`
	}
	SendStatus struct {
		Status       string        `json:"status,omitempty"`
		OrderID      string        `json:"order_id,omitempty"`
		ReceivedTime string        `json:"receivedTime,omitempty"`
		OrderEvents  []OrderEvents `json:"orderEvents"`
	}
	OpenOrder struct {
		OrderID      string  `json:"order_id,omitempty"`
		Symbol       string  `json:"symbol,omitempty"`
		Side         string  `json:"side,omitempty"`
		OrderType    string  `json:"orderType,omitempty"`
		LimitPrice   float64 `json:"limitPrice,omitempty"`
		StopPrice    float64 `json:"stopPrice,omitempty"`
		UnfilledSize float64 `json:"unfilledSize,omitempty"`
		ReceivedTime string  `json:"receivedTime,omitempty"`
		Status       string  `json:"status,omitempty"`
		FilledSize   float64 `json:"filledSize,omitempty"`
	}
	CancelStatus struct {
		Status       string `json:"status,omitempty"`
		OrderID      string `json:"order_id,omitempty"`
		ReceivedTime string `json:"receivedTime,omitempty"`
	}
	ReceiveOrder struct {
		Result        string       `json:"result,omitempty"`
		SendStatus    SendStatus   `json:"sendStatus,omitempty"`
		OpenPositions []OpenOrder  `json:"openOrders,omitempty"`
		CancelStatus  CancelStatus `json:"cancelStatus,omitempty"`
		ServerTime    string       `json:"serverTime,omitempty"`
		Error         string       `json:"error,omitempty"`
	}
)

var ErrOperationNotFound = errors.New("given operation not found")
