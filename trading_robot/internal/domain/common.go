package domain

import (
	"fmt"
	"strconv"
	"time"
)

const (
	UnsubscribePair = "/unsubscribe/{pair}"
	SubscribePair   = "/subscribe/{pair}"
	PairVar         = "pair"

	ShutdownOperation = "/shutdown"
	SetQuantity       = "/quantity/{value}"
	SetMultiplier     = "/multiplier/{value}"
	ValueVal          = "value"
)

type OrderType string

const (
	SellOrder OrderType = "sell"
	BuyOrder  OrderType = "buy"
)

const IocOrder = "ioc"

type UnixTS time.Time

func (t *UnixTS) UnmarshalJSON(data []byte) error {
	i, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	*t = UnixTS(time.UnixMilli(i))
	return nil
}

func (t UnixTS) String() string {
	return fmt.Sprint(time.Time(t))
}

type Price struct {
	Time      UnixTS  `json:"time" validate:"required"`
	ProductID string  `json:"product_id" validate:"required"`
	Quantity  float64 `json:"qty" validate:"required,gte=0"`
	Price     float64 `json:"price" validate:"required,gt=0"`
}

type Order struct {
	OrderID       string  `json:"order_id,omitempty"`
	OrderType     string  `json:"orderType,omitempty"`
	Symbol        string  `json:"symbol,omitempty"`
	Side          string  `json:"side,omitempty"`
	Size          int     `json:"size,omitempty"`
	LimitPrice    float64 `json:"limitPrice,omitempty"`
	StopPrice     float64 `json:"stopPrice,omitempty"`
	TriggerSignal string  `json:"triggerSignal,omitempty"`
	CliOrdID      string  `json:"cliOrdId,omitempty"`
	ReduceOnly    string  `json:"reduceOnly,omitempty"`
}

type CreateOrderResponse struct {
	OrderType      string  `json:"orderType,omitempty"`
	Symbol         string  `json:"symbol,omitempty"`
	Side           string  `json:"side,omitempty"`
	Size           int     `json:"size,omitempty"`
	LimitPrice     float64 `json:"limitPrice,omitempty"`
	Result         string  `json:"result,omitempty"`
	Status         string  `json:"status,omitempty"`
	OrderID        string  `json:"order_id,omitempty"`
	ReceivedTime   string  `json:"receivedTime,omitempty"`
	OrderEventType string  `json:"order_event_type"`
}

func (r CreateOrderResponse) String() string {
	return fmt.Sprintf(`Created new order:
OrderType: %s
Symbol: %s
Side: %s
Size: %v
LimitPrice: %v
Result: %s
Status: %s
OrderID: %s
Time: %s`, r.OrderType, r.Symbol, r.Side, r.Size, r.LimitPrice, r.Result, r.Status, r.OrderID, r.ReceivedTime)
}

func CreateIocOrder(orderType OrderType, pair string, price float64, quantity int) Order {
	return Order{
		OrderType:  IocOrder,
		Symbol:     pair,
		Side:       string(orderType),
		Size:       quantity,
		LimitPrice: price,
	}
}

type Ticker struct {
	Time    UnixTS  `json:"time" validate:"required"`
	Pair    string  `json:"pair" validate:"required"`
	Bid     float64 `json:"bid" validate:"required,gt=0"`
	Ask     float64 `json:"ask" validate:"required,gt=0"`
	BidSize float64 `json:"bid_size" validate:"required,gte=0"`
	AskSize float64 `json:"ask_size" validate:"required,gte=0"`
}
