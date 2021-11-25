package indicator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMACDEvaluator(t *testing.T) {
	a := assert.New(t)

	alphaFunc := func(p int) float64 {
		return 2 / float64(p+1)
	}

	testID := 0
	t.Logf("\tTest %d:\tmacd test", testID)
	{
		macd := NewMACDEvaluator(12, 26, 9, alphaFunc)
		values := []float64{10, 7, 8, 9, 7, 1, 30, 12, 11, 8, 9, 16, 17, 18, 20, 30, 32, 43, 55, 30, 20, 15, 12, 10, 8,
			10, 9, 5, 7, 2, 6, 1, 10, 13, 16, 24, 22, 17, 10, 20, 15, 9, 10, 7, 8, 9, 7, 1, 30, 12, 11, 8, 9, 16, 19,
			20, 15, 16, 17, 18, 20, 30, 32, 43, 55, 30, 20, 15, 12, 10, 8, 10, 9, 5}
		for _, val := range values {
			macd.UpdateMACD(val)
		}
		m, s := macd.GetMACD()
		a.Equalf(-2.0481284566795406, m, "MACD values should be equal")
		a.Equalf(0.7133144817168429, s, "MACD signal should be equal")
	}
}

func TestMACDStrategy(t *testing.T) {
	a := assert.New(t)

	alphaFunc := func(p int) float64 {
		return 2 / float64(p+1)
	}

	testID := 0
	t.Logf("\tTest %d:\tmacd strategy recommends to buy", testID)
	{
		macd := NewMACDEvaluator(12, 26, 9, alphaFunc)
		s := NewMACDStrategy(macd)
		values := []float64{10, 7, 8, 9, 7, 1, 30, 12, 11, 8, 9, 16, 17, 18, 20, 30, 32, 43, 55, 30, 20, 15, 12, 10, 8,
			10, 9, 5, 7, 2, 6, 1, 10, 13, 16}
		for _, val := range values {
			s.Update(val)
		}
		s.Update(24)
		a.Equalf(true, s.Long(), "Indicator should recommend to buy")
		a.Equalf(false, s.Short(), "Indicator should not recommend to sell")
	}

	testID++
	t.Logf("\tTest %d:\tmacd strategy recommends nothing on negative value", testID)
	{
		macd := NewMACDEvaluator(12, 26, 9, alphaFunc)
		s := NewMACDStrategy(macd)
		values := []float64{10, 7, 8, 9, 7, 1, 30, 12, 11, 8, 9, 16, 17, 18, 20, 30, 32, 43, 55, 30, 20, 15, 12, 10, 8,
			10, 9, 5, 7, 2, 6, 1, 10, 13, 16}
		for _, val := range values {
			s.Update(val)
		}
		s.Update(10)
		a.Equalf(false, s.Long(), "Indicator should not recommend to buy")
		a.Equalf(false, s.Short(), "Indicator should not recommend to sell")
	}

	testID++
	t.Logf("\tTest %d:\tmacd strategy recommends nothing on positive value", testID)
	{
		macd := NewMACDEvaluator(12, 26, 9, alphaFunc)
		s := NewMACDStrategy(macd)
		values := []float64{10, 7, 8, 9, 7, 1, 30, 12, 11, 8, 9, 16, 17, 18, 20, 30, 32, 43, 55, 30, 20, 15, 12, 10, 8,
			10, 9, 5, 7, 2, 6, 1, 10, 13, 16, 24, 22, 17, 10, 20, 15, 9}
		for _, val := range values {
			s.Update(val)
		}
		s.Update(15)
		a.Equalf(false, s.Long(), "Indicator should not recommend to buy")
		a.Equalf(false, s.Short(), "Indicator should not recommend to sell")
	}

	testID++
	t.Logf("\tTest %d:\tmacd strategy recommends to sell", testID)
	{
		macd := NewMACDEvaluator(12, 26, 9, alphaFunc)
		s := NewMACDStrategy(macd)
		values := []float64{10, 7, 8, 9, 7, 1, 30, 12, 11, 8, 9, 16, 17, 18, 20, 30, 32, 43, 55, 30, 20, 15, 12, 10, 8,
			10, 9, 5, 7, 2, 6, 1, 10, 13, 16, 24, 22, 17, 10, 20, 15, 9}
		for _, val := range values {
			s.Update(val)
		}
		s.Update(10)
		a.Equalf(false, s.Long(), "Indicator should not recommend to buy")
		a.Equalf(true, s.Short(), "Indicator should recommend to sell")
	}
}
