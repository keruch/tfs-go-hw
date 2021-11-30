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
