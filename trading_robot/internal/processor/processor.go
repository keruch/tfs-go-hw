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
	notifier   OrderNotifier
	logger     *log.Logger

	PriceMultiplier float64 // pointer to change it in runtime
	TradingQuantity int
}

type Repository interface {
	StoreToDB(ctx context.Context, response domain.CreateOrderResponse) error
}

type OrdersSender interface {
	CreateOrder(order domain.Order) (domain.CreateOrderResponse, error)
}

type OrderNotifier interface {
	NotifyUsers(message string)
}

func NewOrdersProcessor(s indicator.Strategy, r Repository, c OrdersSender, n OrderNotifier, l *log.Logger) *OrdersProcessor {
	return &OrdersProcessor{
		strategy:   s,
		repo:       r,
		controller: c,
		notifier:   n,
		logger:     l,

		TradingQuantity: 100,
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
				price *= 1.0 + p.PriceMultiplier
				orderInfo, err = p.controller.CreateOrder(domain.CreateIocOrder(domain.BuyOrder, candle.Ticker, price, p.TradingQuantity))
			} else if p.strategy.Short() {
				price *= 1.0 - p.PriceMultiplier
				orderInfo, err = p.controller.CreateOrder(domain.CreateIocOrder(domain.SellOrder, candle.Ticker, price, p.TradingQuantity))
			}

			if err != nil {
				p.logger.Error(err)
				continue
			}

			if orderInfo.Status == "placed" {
				err = p.repo.StoreToDB(context.Background(), orderInfo)
				if err != nil {
					p.logger.Error(err)
				}
				p.notifier.NotifyUsers(orderInfo.String())
				p.logger.Infof("Created new order: id = %v, price = %v", orderInfo.OrderID, price)
			}
		}
		p.logger.Info("Candles processing done")
	}()

}

func (p *OrdersProcessor) SetPriceMultiplier(m float64) {
	p.PriceMultiplier = m
}

func (p *OrdersProcessor) SetTradingQuantity(q int) {
	p.TradingQuantity = q
}
