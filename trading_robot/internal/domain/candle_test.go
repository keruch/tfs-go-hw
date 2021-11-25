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
		price := Price{
			ProductID: "TEST",
			Price:  100,
		}
		TS := time.Date(2020, time.April, 5, 15, 20, 0, 0, time.Local) // 15:20:00
		candle := NewCandle(price, CandlePeriod1m, TS)
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
		price := Price{
			ProductID: "TEST",
			Price:  200,
		}
		candle = Update(candle, price)
		a.Equalf(200.0, candle.High, "High prices should be equal")
		a.Equalf(200.0, candle.Close, "Close prices should be equal")
	}

	testID++
	t.Logf("\tTest %d:\tUpdate low and close", testID)
	{
		price := Price{
			ProductID: "TEST",
			Price:  50,
		}
		candle = Update(candle, price)
		a.Equalf(50.0, candle.Low, "Low prices should be equal")
		a.Equalf(50.0, candle.Close, "Close prices should be equal")
	}
}

func TestPeriodTS(t *testing.T) {
	a := assert.New(t)

	expectedTime := time.Date(2020, time.April, 5, 15, 20, 0, 0, time.Local) // 15:20:00

	testID := 0
	t.Logf("\tTest %d:\tcandle 1m", testID)
	{
		test1 := expectedTime.Add(20 * time.Second) // 15:20:20
		TS, err := PeriodTS(CandlePeriod1m, test1)
		a.NoErrorf(err, "Should not have error")
		a.Equalf(expectedTime, TS, "Should be the same")
	}

	testID++
	t.Logf("\tTest %d:\tcandle 2m", testID)
	{
		test1 := expectedTime.Add(1 * time.Minute).Add(20 * time.Second) // 15:21:20
		TS, err := PeriodTS(CandlePeriod2m, test1)
		a.NoErrorf(err, "Should not have error")
		a.Equalf(expectedTime, TS, "Should be the same")
	}

	testID++
	t.Logf("\tTest %d:\tcandle 10m", testID)
	{
		test1 := expectedTime.Add(7 * time.Minute).Add(20 * time.Second) // 15:27:20
		TS, err := PeriodTS(CandlePeriod10m, test1)
		a.NoErrorf(err, "Should not have error")
		a.Equalf(expectedTime, TS, "Should be the same")
	}

	testID++
	t.Logf("\tTest %d:\terror invalid period", testID)
	{
		test1 := expectedTime.Add(1 * time.Minute).Add(20 * time.Second) // 15:21:20
		_, err := PeriodTS("20m", test1)
		a.Errorf(err, "Should not have error")
		a.Equalf(ErrUnknownPeriod, err, "Errors should be equal")
	}
}

