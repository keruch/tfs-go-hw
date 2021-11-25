package processor

import (
	"context"
	"sync"

	"github.com/keruch/tfs-go-hw/trading_robot/internal/domain"
	"github.com/keruch/tfs-go-hw/trading_robot/internal/indicator"
	"github.com/keruch/tfs-go-hw/trading_robot/pkg/log"
)

type OrdersProcessor struct {
	strategy   indicator.Strategy
	repo       Repository
	controller OrdersSender
	logger     *log.Logger

	PriceMultiplier *float64 // pointer to change it in runtime
	TradingQuanity  *int
}

type Repository interface {
	StoreToDB(ctx context.Context, response domain.CreateOrderResponse, price float64) error
}

type OrdersSender interface {
	CreateOrder(order domain.Order) (response domain.CreateOrderResponse, err error)
}

func NewOrdersProcessor(s indicator.Strategy, r Repository, c OrdersSender, l *log.Logger) *OrdersProcessor {
	return &OrdersProcessor{
		strategy:   s,
		repo:       r,
		controller: c,
		logger:     l,
	}
}

func (p *OrdersProcessor) ProcessCandles(candles <-chan domain.Candle, wg *sync.WaitGroup) {

	go func() {
		defer wg.Done()
		for candle := range candles {
			p.logger.Trace(candle)

			var (
				orderInfo domain.CreateOrderResponse
				err       error
				price     = candle.Close
			)

			// TODO: add stop-loss/take-profit

			p.strategy.Update(price)

			if p.strategy.Long() {
				price *= 1.0 + *p.PriceMultiplier
				orderInfo, err = p.controller.CreateOrder(domain.CreateIocOrder(domain.BuyOrder, candle.Ticker, price, *p.TradingQuanity))
			} else if p.strategy.Short() {
				price *= 1.0 -* p.PriceMultiplier
				orderInfo, err = p.controller.CreateOrder(domain.CreateIocOrder(domain.SellOrder, candle.Ticker, price, *p.TradingQuanity))
			}

			if err != nil {
				p.logger.Error(err)
				continue
			}

			if orderInfo.Status == "placed" {
				if err = p.repo.StoreToDB(context.Background(), orderInfo, price); err != nil {
					p.logger.Error(err)
					continue
				}
				p.logger.Infof("Created new order: id = %v, price = %v", orderInfo.OrderID, price)
			}
		}
		p.logger.Info("Candles processing done")
	}()

}
