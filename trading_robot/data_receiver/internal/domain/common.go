package domain

type Ticker struct {
	Time      int         `json:"time"`
	Pair      string      `json:"pair"`
	Bid       float64     `json:"bid"`
	Ask       float64     `json:"ask"`
	BidSize   float64     `json:"bid_size"`
	AskSize   float64     `json:"ask_size"`
	OtherData interface{} `json:"-"`
}
