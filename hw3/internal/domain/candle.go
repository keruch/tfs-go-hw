package domain

import (
	"fmt"
	"time"
)

type CandleData interface {
	GetTicker() string
	GetOpen() float64
	GetHigh() float64
	GetLow() float64
	GetClose() float64
	GetTS() time.Time

	// Record function needs to implement handlers.CSVRecord interface
	Record() []string
}

type Candle struct {
	Ticker string
	Period CandlePeriod // Интервал
	Open   float64      // Цена открытия
	High   float64      // Максимальная цена
	Low    float64      // Минимальная цена
	Close  float64      // Цена закрытие
	TS     time.Time    // Время начала интервала
}

func NewCandle(data CandleData, period CandlePeriod) (Candle, error) {
	ts, err := PeriodTS(period, data.GetTS())
	if err != nil {
		return Candle{}, err
	}

	candle := Candle{
		Ticker: data.GetTicker(),
		Period: period,
		Open:   data.GetOpen(),
		High:   data.GetHigh(),
		Low:    data.GetLow(),
		Close:  data.GetClose(),
		TS:     ts,
	}

	return candle, nil
}

func (c *Candle) Update(data CandleData) {
	if c.High < data.GetHigh() {
		c.High = data.GetHigh()
	}

	if c.Low > data.GetLow() {
		c.Low = data.GetLow()
	}

	c.Close = data.GetClose()
}

func (c Candle) String() string {
	TS := c.TS.Format(time.RFC3339)
	return fmt.Sprintf("%s,%s,%f,%f,%f,%f", c.Ticker, TS, c.Open, c.High, c.Low, c.Close)
}

func (c Candle) Record() []string {
	return []string{
		c.Ticker,
		c.TS.Format(time.RFC3339),
		fmt.Sprintf("%f", c.Open),
		fmt.Sprintf("%f", c.High),
		fmt.Sprintf("%f", c.Low),
		fmt.Sprintf("%f", c.Close),
	}
}

func (c Candle) GetTicker() string {
	return c.Ticker
}

func (c Candle) GetTS() time.Time {
	return c.TS
}

func (c Candle) GetOpen() float64 {
	return c.Open
}

func (c Candle) GetHigh() float64 {
	return c.High
}

func (c Candle) GetLow() float64 {
	return c.Low
}

func (c Candle) GetClose() float64 {
	return c.Close
}
