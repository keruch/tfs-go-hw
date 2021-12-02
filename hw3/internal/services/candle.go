package services

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/keruch/tfs-go-hw/hw3/internal/domain"
	"github.com/keruch/tfs-go-hw/hw3/internal/handlers"
	"github.com/keruch/tfs-go-hw/hw3/internal/repository"
	"github.com/keruch/tfs-go-hw/hw3/pkg/log"
)

type CandlesService struct {
	Repo   repository.PricesData
	Logger *log.Logger
}

type tickersMap map[string]domain.Candle

func NewCandlesService(repo repository.PricesData, logger *log.Logger) *CandlesService {
	return &CandlesService{
		Repo:   repo,
		Logger: logger,
	}
}

func (cs *CandlesService) GenerateCandles(ctx context.Context) error {
	in := cs.Repo.GetPrices(ctx)

	var err error
	var wg sync.WaitGroup
	wg.Add(6)

	c1m := cs.generateCandles(in, domain.CandlePeriod1m, &wg)
	c1m, err = cs.writeCandle(c1m, domain.CandlePeriod1m, &wg)
	if err != nil {
		return err
	}

	c2m := cs.generateCandles(c1m, domain.CandlePeriod2m, &wg)
	c2m, err = cs.writeCandle(c2m, domain.CandlePeriod2m, &wg)
	if err != nil {
		return err
	}

	c10m := cs.generateCandles(c2m, domain.CandlePeriod10m, &wg)
	c10m, err = cs.writeCandle(c10m, domain.CandlePeriod10m, &wg)
	if err != nil {
		return err
	}

	// mock function to emulate reading from pipe
	go func() {
		for range c10m {
		}
	}()

	wg.Wait()

	cs.Logger.Info("Generating candles_generator done!")

	return nil
}

func (cs *CandlesService) writeCandle(in <-chan domain.CandleData, period domain.CandlePeriod, wg *sync.WaitGroup) (<-chan domain.CandleData, error) {
	filename := fmt.Sprintf("candles_%s.csv", period)
	csvWriter, err := handlers.NewCSVWriterCloser(filename)
	if err != nil {
		return nil, err
	}

	out := make(chan domain.CandleData)
	go func() {
		defer close(out)
		defer wg.Done()

		for val := range in {
			err = csvWriter.Write(val)
			if err != nil {
				panic("failed to write data to csv")
			}
			csvWriter.Flush()
			out <- val
		}
		cs.Logger.Infof("write done for %s", string(period))
	}()

	return out, nil
}

func (cs *CandlesService) generateCandles(in <-chan domain.CandleData, period domain.CandlePeriod, wg *sync.WaitGroup) <-chan domain.CandleData {
	var (
		out          = make(chan domain.CandleData)
		startPeriod  = true
		tm           tickersMap
		prevPeriodTS time.Time
	)

	go func() {
		defer close(out)
		defer wg.Done()

		for val := range in {
			candle, err := domain.NewCandle(val, period)
			if err != nil {
				panic("failed to create new candle")
			}

			if candle.TS != prevPeriodTS {
				prevPeriodTS = candle.TS
				startPeriod = true
			}

			if startPeriod {
				flushDataToChan(tm, out)
				tm = make(tickersMap) // clear map
				startPeriod = false
			}

			tm = updateTickersMap(tm, candle)
		}
		flushDataToChan(tm, out)
		cs.Logger.Infof("generate done for %s", string(period))
	}()

	return out
}

func flushDataToChan(tm tickersMap, out chan domain.CandleData) {
	for _, mapVal := range tm {
		out <- mapVal
	}
}

func updateTickersMap(tm tickersMap, candle domain.Candle) tickersMap {
	tickerVal, ok := tm[candle.GetTicker()]
	if !ok {
		tickerVal = candle
	} else {
		tickerVal.Update(candle)
	}

	tm[candle.GetTicker()] = tickerVal
	return tm
}
