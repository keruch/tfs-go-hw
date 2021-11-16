package domain

type SendOrder struct {
	OrderID       string `json:"order_id,omitempty"`
	OrderType     string `json:"orderType,omitempty"`
	Symbol        string `json:"symbol,omitempty"`
	Side          string `json:"side,omitempty"`
	Size          string `json:"size,omitempty"`
	LimitPrice    string `json:"limitPrice,omitempty"`
	StopPrice     string `json:"stopPrice,omitempty"`
	TriggerSignal string `json:"triggerSignal,omitempty"`
	CliOrdId      string `json:"cliOrdId,omitempty"`
	ReduceOnly    string `json:"reduceOnly,omitempty"`
}

type SendStatus struct {
	Status       string                 `json:"status,omitempty"`
	OrderID      string                 `json:"order_id,omitempty"`
	ReceivedTime string                 `json:"receivedTime,omitempty"`
	OtherData    map[string]interface{} `json:"-"`
}

type OpenOrder struct {
	OrderID      string                 `json:"order_id,omitempty"`
	Symbol       string                 `json:"symbol,omitempty"`
	Side         string                 `json:"side,omitempty"`
	OrderType    string                 `json:"orderType,omitempty"`
	LimitPrice   float64                `json:"limitPrice,omitempty"`
	StopPrice    float64                `json:"stopPrice,omitempty"`
	UnfilledSize float64                `json:"unfilledSize,omitempty"`
	ReceivedTime string                 `json:"receivedTime,omitempty"`
	Status       string                 `json:"status,omitempty"`
	FilledSize   float64                `json:"filledSize,omitempty"`
	OtherData    map[string]interface{} `json:"-"`
}

type EditStatus struct {
	Status       string                 `json:"status,omitempty"`
	OrderID      string                 `json:"order_id,omitempty"`
	ReceivedTime string                 `json:"receivedTime,omitempty"`
	OtherData    map[string]interface{} `json:"-"`
}

type CancelStatus struct {
	Status       string                 `json:"status,omitempty"`
	OrderID      string                 `json:"order_id,omitempty"`
	ReceivedTime string                 `json:"receivedTime,omitempty"`
	OtherData    map[string]interface{} `json:"-"`
}

type ReceiveOrder struct {
	Result        string       `json:"result,omitempty"`
	SendStatus    SendStatus   `json:"sendStatus,omitempty"`
	OpenPositions []OpenOrder  `json:"openOrders,omitempty"`
	EditStatus    EditStatus   `json:"editStatus,omitempty"`
	CancelStatus  CancelStatus `json:"cancelStatus,omitempty"`
	ServerTime    string       `json:"serverTime,omitempty"`
	Error         string       `json:"error,omitempty"`
}
