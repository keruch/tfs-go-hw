package candles

import (
	"time"

	"github.com/keruch/tfs-go-hw/trading_robot/internal/domain"
)

func GenerateCandles(in <-chan domain.Ticker, period domain.CandlePeriod) <-chan domain.Candle {
	out := make(chan domain.Candle)

	go func() {
		defer close(out)

		var (
			startPeriod = false
			currentTS   time.Time
			candle      domain.Candle
		)
		for ticker := range in {
			candleTS, err := domain.PeriodTS(period, ticker.Time)
			if err != nil {
				panic(err)
			}

			if candleTS != currentTS {
				currentTS = candleTS
				startPeriod = true
			}

			if startPeriod {
				out <- candle
				candle = domain.NewCandle(ticker, period, candleTS)
				startPeriod = false
			}

			candle = domain.Update(candle, ticker)
		}
		out <- candle
	}()

	return out
}
