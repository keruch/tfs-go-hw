package candles

import (
	"testing"
	"time"

	"github.com/keruch/tfs-go-hw/trading_robot/internal/domain"
	"github.com/stretchr/testify/assert"
)

var (
	mockTime = time.Date(2020, time.April, 5, 15, 20, 10, 0, time.Local) // 15:20:10
	mockPair = "TEST_TICKER"
)

var (
	test1 = domain.Candle{
		Ticker: mockPair,
		Period: domain.CandlePeriod1m,
		Open:   100,
		High:   120,
		Low:    90,
		Close:  110,
		TS:     time.Date(2020, time.April, 5, 15, 20, 0, 0, time.Local), // 15:20:00
	}

	test2 = domain.Candle{
		Ticker: mockPair,
		Period: domain.CandlePeriod1m,
		Open:   100,
		High:   100,
		Low:    90,
		Close:  90,
		TS:     time.Date(2020, time.April, 5, 15, 22, 0, 0, time.Local), // 15:22:00
	}

	test3 = domain.Candle{
		Ticker: mockPair,
		Period: domain.CandlePeriod1m,
		Open:   100,
		High:   100,
		Low:    100,
		Close:  100,
		TS:     time.Date(2020, time.April, 5, 15, 24, 0, 0, time.Local), // 15:24:00,
	}
)

func TestGenerateCandles(t *testing.T) {
	a := assert.New(t)

	in := MockTickersGenerator()
	out := GenerateCandles(in, domain.CandlePeriod1m)

	_ = <-out // the first value is always invalid

	testID := 0
	t.Logf("\tTest %d:\tOne full candle", testID)
	{
		candle := <-out
		a.Equalf(test1, candle, "Tickers should be the same")
	}

	testID++
	t.Logf("\tTest %d:\tOpen equal high, low equal close", testID)
	{
		candle := <-out
		a.Equalf(test2, candle, "Tickers should be the same")
	}

	testID++
	t.Logf("\tTest %d:\tOpen equal high, low and close", testID)
	{
		candle := <-out
		a.Equalf(test3, candle, "Tickers should be the same")
	}
}

func MockTickersGenerator() <-chan domain.Ticker {
	out := make(chan domain.Ticker)

	tickers := []domain.Ticker{
		// for test 1
		{
			Time: mockTime.Add(10 * time.Second), // 15:20:20
			Pair: mockPair,
			Ask:  test1.Open,
		},
		{
			Time: mockTime.Add(20 * time.Second), // 15:20:30,
			Pair: mockPair,
			Ask:  test1.High,
		},
		{
			Time: mockTime.Add(30 * time.Second), // 15:20:40,
			Pair: mockPair,
			Ask:  test1.Low,
		},
		{
			Time: mockTime.Add(40 * time.Second), // 15:20:50,
			Pair: mockPair,
			Ask:  test1.Close,
		},
		// tickers for test 2
		{
			Time: mockTime.Add(2 * time.Minute).Add(10 * time.Second), // 15:22:20,
			Pair: mockPair,
			Ask:  test2.Open,
		},
		{
			Time: mockTime.Add(2 * time.Minute).Add(20 * time.Second), // 15:22:30,
			Pair: mockPair,
			Ask:  test2.Close,
		},
		// tickers for test 3
		{
			Time: mockTime.Add(4 * time.Minute).Add(40 * time.Second), // 15:24:50,
			Pair: "TEST_TICKER",
			Ask:  test3.Open,
		},
	}

	go func() {
		defer close(out)
		for _, ticker := range tickers {
			out <- ticker
		}
	}()

	return out
}
