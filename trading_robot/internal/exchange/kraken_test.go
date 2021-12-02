package exchange

import (
	"testing"

	"github.com/keruch/tfs-go-hw/trading_robot/pkg/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type krakenEnvironment struct {
	suite.Suite
	ex *KrakenExchange
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

func (k *krakenEnvironment) TestSubscribePairs() {
	testID := 0
	k.T().Logf("\tTest %d:\tsubscribe pairs success", testID)
	{
		err := k.ex.SubscribePairs("TEST_PAIR")
		k.NoError(err)
		k.Equal(1, len(k.ex.pairs))
	}

	testID++
	k.T().Logf("\tTest %d:\tsubscribe pairs nore than one", testID)
	{
		err := k.ex.SubscribePairs("TEST_PAIR")
		k.Equal(1, len(k.ex.pairs))
		k.Error(err)
	}

	testID++
	k.T().Logf("\tTest %d:\tunsubscribe pairs success", testID)
	{
		err := k.ex.UnsubscribePairs("TEST_PAIR")
		k.NoError(err)
		k.Equal(0, len(k.ex.pairs))
	}
}

func (k *krakenEnvironment) TearDownSuite() {
	k.ex.CloseConnection()
}

func TestDatabase(t *testing.T) {
	suite.Run(t, new(krakenEnvironment))
}
