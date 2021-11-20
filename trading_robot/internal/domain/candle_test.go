package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewCandle(t *testing.T) {
	a := assert.New(t)

	testID := 0
	t.Logf("\tTest %d:\tCreate new candle", testID)
	{
		ticker := Ticker{
			Pair: "TEST",
			Ask:  100,
		}
		TS := time.Date(2020, time.April, 5, 15, 20, 0, 0, time.Local) // 15:20:00
		candle := NewCandle(ticker, CandlePeriod1m, TS)
		expectedCandle := Candle{
			Ticker: "TEST",
			Period: CandlePeriod1m,
			Open:   100,
			High:   100,
			Low:    100,
			Close:  100,
			TS:     TS,
		}
		a.Equalf(expectedCandle, candle, "Candles should be equal")
	}
}

func TestUpdate(t *testing.T) {
	a := assert.New(t)

	TS := time.Date(2020, time.April, 5, 15, 20, 0, 0, time.Local) // 15:20:00
	candle := Candle{
		Ticker: "TEST",
		Period: CandlePeriod1m,
		Open:   100,
		High:   100,
		Low:    100,
		Close:  100,
		TS:     TS,
	}

	testID := 0
	t.Logf("\tTest %d:\tUpdate high and close", testID)
	{
		ticker := Ticker{
			Pair: "TEST",
			Ask:  200,
		}
		candle = Update(candle, ticker)
		a.Equalf(200.0, candle.High, "High prices should be equal")
		a.Equalf(200.0, candle.Close, "Close prices should be equal")
	}

	testID++
	t.Logf("\tTest %d:\tUpdate low and close", testID)
	{
		ticker := Ticker{
			Pair: "TEST",
			Ask:  50,
		}
		candle = Update(candle, ticker)
		a.Equalf(50.0, candle.Low, "Low prices should be equal")
		a.Equalf(50.0, candle.Close, "Close prices should be equal")
	}
}
