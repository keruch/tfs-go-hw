package domain

const (
	UnsubscribePair = "/unsubscribe/{pair}"
	SubscribePair   = "/subscribe/{pair}"
	PairVar         = "pair"

	IndexEndpoint      = "/index/{params}"
	OperationsEndpoint = "/operations"

	ShutdownOperation = "/shutdown"
	StartOperation    = "/start"
	StopOperation     = "/stop"
)
