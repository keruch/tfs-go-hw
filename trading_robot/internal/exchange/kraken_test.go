package exchange

import (
	"testing"

	"github.com/keruch/tfs-go-hw/trading_robot/internal/exchange/kraken"
	"github.com/keruch/tfs-go-hw/trading_robot/pkg/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type krakenEnvironment struct {
	suite.Suite
	ex Exchange
}

func (k *krakenEnvironment) SetupSuite() {
	logger := &log.Logger{}
	k.ex, _ = NewKrakenExchange(logger)
}



func TestNewKrakenExchange(t *testing.T) {
	a := assert.New(t)

	testID := 0
	t.Logf("\tTest %d:\tcreate kraken exchange", testID)
	{
		logger := &log.Logger{}
		_, err := NewKrakenExchange(logger)
		a.NoErrorf(err, "Should create without error")
	}
}

func TestQueryByOperation(t *testing.T) {
	a := assert.New(t)

	testID := 0
	t.Logf("\tTest %d:\tcreate order operation", testID)
	{
		_, err := QueryByOperation(&kraken.SendOrder{}, kraken.CreateOrder)
		a.NoErrorf(err, "Should not be error")
	}
	testID++
	t.Logf("\tTest %d:\topen orders operation", testID)
	{
		_, err := QueryByOperation(&kraken.SendOrder{}, kraken.OpenOrders)
		a.NoErrorf(err, "Should not be error")
	}
	testID++
	t.Logf("\tTest %d:\tcancel order operation", testID)
	{
		_, err := QueryByOperation(&kraken.SendOrder{}, kraken.CreateOrder)
		a.NoErrorf(err, "Should not be error")
	}
	testID++
	t.Logf("\tTest %d:\tcancel order operation", testID)
	{
		_, err := QueryByOperation(&kraken.SendOrder{}, "TEST_OPERATION")
		a.Errorf(err, "Should not error")
		a.Equalf(kraken.ErrOperationNotFound, err, "Errors should be equal")
	}
}
