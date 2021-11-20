package kraken

import (
	"errors"

	"github.com/keruch/tfs-go-hw/trading_robot/internal/domain"
)

const (
	MaxRetries = 13

	Scheme = "https"
	Host   = "demo-futures.kraken.com"
	Path   = "/derivatives"

	WsScheme = "wss"
	WsPath   = "/ws/v1"
	WsQuery  = "chart"
)

type (
	Event string
	Feed  string
)

const (
	SubscribeEvent   Event = "subscribe"
	UnsubscribeEvent Event = "unsubscribe"
	Candles1mFeed    Feed  = "candles_trade_1m"
	TickerFeed       Feed  = "ticker"
)

type OperationEndpoint string

const (
	CreateOrder OperationEndpoint = "/api/v3/sendorder"
	OpenOrders  OperationEndpoint = "/api/v3/openorders"
	EditOrder   OperationEndpoint = "/api/v3/editorder" // does not supported because of orderId filed
	CancelOrder OperationEndpoint = "/api/v3/cancelorder"
)

type RequestHeader = string

const (
	Authent RequestHeader = "Authent"
	APIKey  RequestHeader = "APIKey"
)

type CandleData struct {
	C domain.Candle `json:"candle"`
}

var ErrOperationNotFound = errors.New("given operation not found")
