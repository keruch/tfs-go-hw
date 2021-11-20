package domain

import "time"

type Ticker struct {
	Time      time.Time   `json:"time"`
	Pair      string      `json:"pair"`
	Bid       float64     `json:"bid"`
	Ask       float64     `json:"ask"`
	BidSize   float64     `json:"bid_size"`
	AskSize   float64     `json:"ask_size"`
	OtherData interface{} `json:"-"`
}

const (
	WaitTime time.Duration = time.Second * 5
)
