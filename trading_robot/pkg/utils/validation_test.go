package utils

import (
	"testing"
	"time"

	"github.com/keruch/tfs-go-hw/trading_robot/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestValidateDataIsPrice(t *testing.T) {
	a := assert.New(t)

	testID := 0
	t.Logf("\tTest %d:\tsubscription price data", testID)
	{
		price, ok := ValidateDataIsPrice([]byte(subscriptionPriceData))
		a.Equalf(true, ok, "Should subscribe")
		expectedPrice := domain.Price{
			Time:      domain.UnixTS(time.UnixMilli(1612266317519)),
			ProductID: "PI_XBTUSD",
			Quantity:  15000,
			Price:     34969.5,
		}
		a.Equalf(expectedPrice, price, "Should be equal")
	}

	testID++
	t.Logf("\tTest %d:\tsubscribed event", testID)
	{
		price, ok := ValidateDataIsPrice([]byte(subscribedEvent))
		a.Equalf(false, ok, "Subscribed event should not be valid")
		a.Equalf(domain.Price{}, price, "Ticker should be nil")
	}

	testID++
	t.Logf("\tTest %d:\tWS connected event", testID)
	{
		price, ok := ValidateDataIsPrice([]byte(connectedWS))
		a.Equalf(false, ok, "WS connected event should not be valid")
		a.Equalf(domain.Price{}, price, "Ticker should be nil")
	}

	testID++
	t.Logf("\tTest %d:\terror event", testID)
	{
		price, ok := ValidateDataIsPrice([]byte(connectedWS))
		a.Equalf(false, ok, "Error event should not be valid")
		a.Equalf(domain.Price{}, price, "Ticker should be nil")
	}
}

func TestValidateDataIsTicker(t *testing.T) {
	a := assert.New(t)

	testID := 0
	t.Logf("\tTest %d:\tsubscription ticker data", testID)
	{
		ticker, ok := ValidateDataIsTicker([]byte(subscriptionTickerData))
		a.Equalf(true, ok, "Should subscribe")
		expectedTicker := domain.Ticker{
			Time:    domain.UnixTS(time.UnixMilli(1612270825253)),
			Pair:    "XBT:USD",
			Bid:     34832.5,
			Ask:     34847.5,
			BidSize: 42864,
			AskSize: 2300,
		}
		a.Equalf(expectedTicker, ticker, "Should be equal")
	}

	testID++
	t.Logf("\tTest %d:\tsubscribed event", testID)
	{
		ticker, ok := ValidateDataIsTicker([]byte(subscribedEvent))
		a.Equalf(false, ok, "Subscribed event should not be valid")
		a.Equalf(domain.Ticker{}, ticker, "Ticker should be nil")
	}

	testID++
	t.Logf("\tTest %d:\tWS connected event", testID)
	{
		ticker, ok := ValidateDataIsTicker([]byte(connectedWS))
		a.Equalf(false, ok, "WS connected event should not be valid")
		a.Equalf(domain.Ticker{}, ticker, "Ticker should be nil")
	}

	testID++
	t.Logf("\tTest %d:\terror event", testID)
	{
		ticker, ok := ValidateDataIsTicker([]byte(connectedWS))
		a.Equalf(false, ok, "Error event should not be valid")
		a.Equalf(domain.Ticker{}, ticker, "Ticker should be nil")
	}
}

const (
	subscribedEvent = `{  
    "event":"subscribed",
    "feed":"ticker",
    "product_ids":[  
        "PI_XBTUSD"
    ]
}`

	connectedWS = `{
    "event": "info",
    "version": 1
}`

	subscriptionTickerData = `{
  "time": 1612270825253,
  "feed": "ticker",
  "product_id": "PI_XBTUSD",
  "bid": 34832.5,
  "ask": 34847.5,
  "bid_size": 42864,
  "ask_size": 2300,
  "volume": 262306237,
  "dtm": 0,
  "leverage": "50x",
  "indicator": 34803.45,
  "premium": 0.1,
  "last": 34852,
  "change": 2.995109121267192,
  "funding_rate": 3.891007752e-9,
  "funding_rate_prediction": 4.2233756e-9,
  "suspended": false,
  "tag": "perpetual",
  "pair": "XBT:USD",
  "openInterest": 107706940,
  "markPrice": 34844.25,
  "maturityTime": 0,
  "relative_funding_rate": 0.000135046879166667,
  "relative_funding_rate_prediction": 0.000146960125,
  "next_funding_rate_time": 1612281600000
}`

	subscriptionPriceData = `{
  "feed": "trade",
  "product_id": "PI_XBTUSD",
  "uid": "05af78ac-a774-478c-a50c-8b9c234e071e",
  "side": "sell",
  "type": "fill",
  "seq": 653355,
  "time": 1612266317519,
  "qty": 15000,
  "price": 34969.5
}`
)
