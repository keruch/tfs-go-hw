package processor

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/keruch/tfs-go-hw/trading_robot/config"
	"github.com/keruch/tfs-go-hw/trading_robot/internal/domain"
	"github.com/keruch/tfs-go-hw/trading_robot/pkg/log"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type RepoMock struct {
	mock.Mock
}

func (r *RepoMock) StoreToDB(ctx context.Context, response domain.CreateOrderResponse) error {
	args := r.Called(ctx, response)
	return args.Error(0)
}

type OrdersSenderMock struct {
	mock.Mock
}

func (c *OrdersSenderMock) CreateOrder(order domain.Order) (response domain.CreateOrderResponse, err error) {
	args := c.Called(order)
	return args.Get(0).(domain.CreateOrderResponse), args.Error(1)
}

type StrategyMock struct {
	mock.Mock
}

func (s *StrategyMock) Update(p float64) {
	s.Called(p)
}

func (s *StrategyMock) Long() bool {
	args := s.Called()
	return args.Bool(0)
}

func (s *StrategyMock) Short() bool {
	args := s.Called()
	return args.Bool(0)
}

type NotifierMock struct {
	mock.Mock
}

func (n *NotifierMock) NotifyUsers(message string) {
	n.Called(message)
}

type Environment struct {
	suite.Suite
	repo       *RepoMock
	controller *OrdersSenderMock
	strategy   *StrategyMock
	notifier   *NotifierMock
}

func (e *Environment) SetupSuite() {
	if err := config.SetupConfig(); err != nil {
		e.T().Fatal(err)
	}

	e.controller = new(OrdersSenderMock)
	e.repo = new(RepoMock)
	e.strategy = new(StrategyMock)
	e.notifier = new(NotifierMock)
}

func (e *Environment) TearDownSuite() {
	e.controller.AssertExpectations(e.T())
	e.repo.AssertExpectations(e.T())
	e.strategy.AssertExpectations(e.T())
	e.notifier.AssertExpectations(e.T())
}

var validResponse = domain.CreateOrderResponse{
	OrderType:    "ioc",
	Symbol:       "TEST_SYM",
	Side:         "sell",
	Size:         100,
	LimitPrice:   4213.1,
	Result:       "success",
	Status:       "placed",
	OrderID:      "8dcdbe17-b729-4fef-8b89-36e561535f38",
	ReceivedTime: "2021-11-25T19:05:03.670Z",
}

func (e *Environment) TestProcessor() {
	logger := log.NewLogger()
	logger.SetLevel(0) // set panic level to prevent output spam
	processor := NewOrdersProcessor(e.strategy, e.repo, e.controller, e.notifier, logger)

	testID := 0
	e.T().Logf("\tTest %d:\tprocessor all success long", testID)
	{
		e.strategy.On("Update", mock.Anything).Once()
		e.strategy.On("Long").Return(true).Once()
		e.controller.On("CreateOrder", mock.Anything).Return(validResponse, nil).Once()
		e.repo.On("StoreToDB", mock.Anything, mock.Anything).Return(nil).Once()
		e.notifier.On("NotifyUsers", mock.Anything).Return().Once()
	}

	testID++
	e.T().Logf("\tTest %d:\tprocessor all success short", testID)
	{
		e.strategy.On("Update", mock.Anything).Once()
		e.strategy.On("Long").Return(false).Once()
		e.strategy.On("Short").Return(true).Once()
		e.controller.On("CreateOrder", mock.Anything).Return(validResponse, nil).Once()
		e.repo.On("StoreToDB", mock.Anything, mock.Anything).Return(nil).Once()
		e.notifier.On("NotifyUsers", mock.Anything).Return().Once()
	}

	testID++
	e.T().Logf("\tTest %d:\tprocessor createIocOrder failed", testID)
	{
		e.strategy.On("Update", mock.Anything).Once()
		e.strategy.On("Long").Return(true).Once()
		e.controller.On("CreateOrder", mock.Anything).Return(domain.CreateOrderResponse{}, errors.New("CreateOrder failed")).Once()
	}

	testID++
	e.T().Logf("\tTest %d:\tprocessor StoreToDB failed", testID)
	{
		e.strategy.On("Update", mock.Anything).Once()
		e.strategy.On("Long").Return(true).Once()
		e.controller.On("CreateOrder", mock.Anything).Return(validResponse, nil).Once()
		e.repo.On("StoreToDB", mock.Anything, mock.Anything).Return(errors.New("store error")).Once()
		e.notifier.On("NotifyUsers", mock.Anything).Return().Once()
	}

	candles := []domain.Candle{{Close: 4, Ticker: "TEST"}, {Close: 5, Ticker: "TEST"}, {Close: 8, Ticker: "TEST"}, {Close: 10, Ticker: "TEST"}}
	out := make(chan domain.Candle)
	go func() {
		defer close(out)
		for _, candle := range candles {
			out <- candle
		}
	}()
	var wg sync.WaitGroup
	wg.Add(1)
	processor.ProcessCandles(out, &wg)
	wg.Wait()
}

func TestOrdersProcessor_ProcessCandles(t *testing.T) {
	suite.Run(t, new(Environment))
}
