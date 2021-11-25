package domain

import (
	"errors"
	"time"
)

// Candle TODO: move to candles path
type Candle struct {
	Ticker string
	Period CandlePeriod // Интервал
	Open   float64      // Цена открытия
	High   float64      // Максимальная цена
	Low    float64      // Минимальная цена
	Close  float64      // Цена закрытие
	TS     time.Time    // Время начала интервала
}

func NewCandle(price Price, period CandlePeriod, TS time.Time) Candle {
	return Candle{
		Ticker: price.ProductID,
		Period: period,
		Open: price.Price,
		High: price.Price,
		Low: price.Price,
		Close: price.Price,
		TS:     TS,
	}
}

func Update(c Candle, p Price) Candle {
	if c.High < p.Price {
		c.High = p.Price
	}

	if c.Low > p.Price || c.Low == 0 {
		c.Low = p.Price
	}

	c.Close = p.Price

	return c
}

var ErrUnknownPeriod = errors.New("unknown period")

type CandlePeriod string

const (
	CandlePeriod1m  CandlePeriod = "1m"
	CandlePeriod2m  CandlePeriod = "2m"
	CandlePeriod10m CandlePeriod = "10m"
)

func PeriodTS(period CandlePeriod, ts time.Time) (time.Time, error) {
	switch period {
	case CandlePeriod1m:
		return ts.Truncate(time.Minute), nil
	case CandlePeriod2m:
		return ts.Truncate(2 * time.Minute), nil
	case CandlePeriod10m:
		return ts.Truncate(10 * time.Minute), nil
	default:
		return time.Time{}, ErrUnknownPeriod
	}
}
