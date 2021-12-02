package indicator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEMA(t *testing.T) {
	a := assert.New(t)

	testID := 0
	t.Logf("\tTest %d:\tema test", testID)
	{
		period := 7
		alpha := 2 / float64(period+1)
		ema := EMA(4, 15, alpha)
		a.Equalf(6.75, ema, "Should be equal")
	}
}

func TestEMAEvaluator(t *testing.T) {
	a := assert.New(t)

	alphaFunc := func(p int) float64 {
		return 2 / float64(p+1)
	}

	testID := 0
	t.Logf("\tTest %d:\tema random values", testID)
	{
		period := 6
		e := NewEMAEvaluator(period, alphaFunc)
		values := []float64{4, 2, 3, 5, 7, 15}
		for _, val := range values {
			e.UpdateEMA(val)
		}
		a.Equalf(7.648003807937169, e.GetEMA(), "Should be equal")
	}

	testID++
	t.Logf("\tTest %d:\tema equal values", testID)
	{
		period := 8
		alphaFunc := func(p int) float64 {
			return 2 / float64(p+1)
		}
		e := NewEMAEvaluator(period, alphaFunc)
		values := []float64{3, 3, 3, 3, 3}
		for _, val := range values {
			e.UpdateEMA(val)
		}
		a.Equalf(3.0, e.GetEMA(), "Should be equal")
	}

	testID++
	t.Logf("\tTest %d:\tema mean values", testID)
	{
		period := 4
		alphaFunc := func(p int) float64 {
			return 2 / float64(p+1)
		}
		e := NewEMAEvaluator(period, alphaFunc)
		values := []float64{2, 12, 2, 12, 2, 12}
		for _, val := range values {
			e.UpdateEMA(val)
		}
		a.Equalf(7.958400000000001, e.GetEMA(), "Should be equal")
	}
}

func TestEMAStrategy(t *testing.T) {
	a := assert.New(t)

	period := 4
	alphaFunc := func(p int) float64 {
		return 2 / float64(p+1)
	}

	testID := 0
	t.Logf("\tTest %d:\tema strategy buy", testID)
	{
		e := NewEMAEvaluator(period, alphaFunc)
		s := NewEMAStrategy(e)
		values := []float64{4, 2, 3, 5, 7, 15}
		for _, val := range values {
			s.Update(val)
		}
		s.Update(15)
		a.Equalf(true, s.Long(), "Indicator should recommend to buy")
		a.Equalf(false, s.Short(), "Indicator should not recommend to sell")
	}

	testID++
	t.Logf("\tTest %d:\tema buy", testID)
	{
		e := NewEMAEvaluator(period, alphaFunc)
		s := NewEMAStrategy(e)
		values := []float64{4, 2, 3, 5, 7, 15}
		for _, val := range values {
			s.Update(val)
		}
		s.Update(5)
		a.Equalf(false, s.Long(), "Indicator should not recommend to buy")
		a.Equalf(true, s.Short(), "Indicator should recommend to sell")
	}
}
