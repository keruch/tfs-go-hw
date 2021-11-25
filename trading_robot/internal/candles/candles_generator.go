package candles

import (
	"sync"
	"time"

	"github.com/keruch/tfs-go-hw/trading_robot/internal/domain"
)

// GenerateCandles TODO: fix invalid first value
func GenerateCandles(in <-chan domain.Price, period domain.CandlePeriod, wg *sync.WaitGroup) <-chan domain.Candle {
	out := make(chan domain.Candle)

	go func() {
		defer wg.Done()
		defer close(out)

		var (
			startPeriod = false
			currentTS   time.Time
			candle      domain.Candle
		)
		for price := range in {
			candleTS, err := domain.PeriodTS(period, time.Time(price.Time))
			if err != nil {
				panic(err)
			}

			if candleTS != currentTS {
				currentTS = candleTS
				startPeriod = true
			}

			if startPeriod {
				candle = domain.NewCandle(price, period, candleTS)
				out <- candle
				startPeriod = false
			}

			candle = domain.Update(candle, price)
		}
		out <- candle
	}()

	return out
}
