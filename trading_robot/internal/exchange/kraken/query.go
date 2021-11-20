package kraken

type QueryParam = string
type QueryParams map[QueryParam]string

// strings for send order
const (
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
