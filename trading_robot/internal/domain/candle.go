package domain

import (
	"time"
)

type Candle struct {
	Ticker string
	Period CandlePeriod // Интервал
	Open   float64      // Цена открытия
	High   float64      // Максимальная цена
	Low    float64      // Минимальная цена
	Close  float64      // Цена закрытие
	TS     time.Time    // Время начала интервала
}

func NewCandle(ticker Ticker, period CandlePeriod, TS time.Time) Candle {
	return Candle{
		Ticker: ticker.Pair,
		Period: period,
		Open: ticker.Ask,
		High: ticker.Ask,
		Low: ticker.Ask,
		Close: ticker.Ask,
		TS:     TS,
	}
}

func Update(c Candle, ticker Ticker) Candle {
	if c.High < ticker.Ask {
		c.High = ticker.Ask
	}

	if c.Low > ticker.Ask || c.Low == 0 {
		c.Low = ticker.Ask
	}

	c.Close = ticker.Ask

	return c
}

// TODO: uncomment if needed
//func (c Candle) String() string {
//	TS := c.TS.Format(time.RFC3339)
//	return fmt.Sprintf("%s,%s,%f,%f,%f,%f", c.Ticker, TS, c.Open, c.High, c.Low, c.Close)
//}
