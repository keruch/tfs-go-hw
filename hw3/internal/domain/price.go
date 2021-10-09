package domain

import (
	"fmt"
	"time"
)

type Price struct {
	Ticker string
	Value  float64
	TS     time.Time
}

func (p Price) String() string {
	return fmt.Sprintf("Ticker:%s Value:%v TS:%s", p.Ticker, p.Value, p.TS)
}

func (p Price) Record() []string {
	return []string{
		p.Ticker,
		fmt.Sprintf("%f", p.Value),
		p.TS.Format(time.RFC3339),
	}
}

func (p Price) GetTicker() string {
	return p.Ticker
}

func (p Price) GetTS() time.Time {
	return p.TS
}

func (p Price) GetOpen() float64 {
	return p.Value
}

func (p Price) GetHigh() float64 {
	return p.Value
}

func (p Price) GetLow() float64 {
	return p.Value
}

func (p Price) GetClose() float64 {
	return p.Value
}
