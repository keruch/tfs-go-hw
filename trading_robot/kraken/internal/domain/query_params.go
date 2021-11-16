package domain

type QueryParams map[string]string

// strings for send order
const (
	OrderID       = "order_id"
	OrderType     = "orderType"
	Symbol        = "symbol"
	Side          = "side"
	Size          = "size"
	LimitPrice    = "limitPrice"
	StopPrice     = "stopPrice"
	TriggerSignal = "triggerSignal"
	CliOrdId      = "cliOrdId"
	ReduceOnly    = "reduceOnly"
)
