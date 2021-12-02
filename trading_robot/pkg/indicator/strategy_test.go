package indicator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMacdEmaStrategy(t *testing.T) {
	a := assert.New(t)

	alphaFunc := func(p int) float64 {
		return 2 / float64(p+1)
	}

	testID := 0
	t.Logf("\tTest %d:\tstrategy composition recommends to sell", testID)
	{
		period := 100
		me := NewMACDEvaluator(12, 26, 9, alphaFunc)
		ee := NewEMAEvaluator(period, alphaFunc)
		sm := NewMACDStrategy(me)
		se := NewEMAStrategy(ee)
		sc := NewStrategiesComposition(se, sm)
		values := []float64{10, 7, 8, 9, 7, 1, 30, 12, 11, 8, 9, 16, 17, 18, 20, 30, 32, 43, 55, 30, 20, 15, 12, 10, 8,
			10, 9, 5, 7, 2, 6, 1, 10, 13, 16, 24, 22, 17, 10, 20, 15, 9}
		for _, val := range values {
			sc.Update(val)
		}
		sc.Update(10)
		a.Equalf(false, sc.Long(), "Strategy composition should not recommend to buy")
		a.Equalf(true, sc.Short(), "Strategy composition should recommend to sell")
	}

	testID++
	t.Logf("\tTest %d:\tstrategy composition recommends nothing (one - long, other - short)", testID)
	{
		period := 100
		me := NewMACDEvaluator(12, 26, 9, alphaFunc)
		ee := NewEMAEvaluator(period, alphaFunc)
		sm := NewMACDStrategy(me)
		se := NewEMAStrategy(ee)
		sc := NewStrategiesComposition(se, sm)
		values := []float64{10, 7, 8, 9, 7, 1, 30, 12, 11, 8, 9, 16, 17, 18, 20, 30, 32, 43, 55, 30, 20, 15, 12, 10, 8,
			10, 9, 5, 7, 2, 6, 1, 10, 13, 16, 24, 22, 17, 10, 20, 15, 9}
		for _, val := range values {
			sc.Update(val)
		}
		sc.Update(15)
		a.Equalf(false, sc.Long(), "Indicator should not recommend to buy")
		a.Equalf(false, sc.Short(), "Indicator should not recommend to sell")
	}

	testID++
	t.Logf("\tTest %d:\tstrategy composition recommends nothing (one - long, other - short)", testID)
	{
		period := 100
		me := NewMACDEvaluator(12, 26, 9, alphaFunc)
		ee := NewEMAEvaluator(period, alphaFunc)
		sm := NewMACDStrategy(me)
		se := NewEMAStrategy(ee)
		sc := NewStrategiesComposition(se, sm)
		values := []float64{10, 7, 8, 9, 7, 1, 30, 12, 11, 8, 9, 16, 17, 18, 20, 30, 32, 43, 55, 30, 20, 15, 12, 10, 8,
			10, 9, 5, 7, 2, 6, 1, 10, 13, 16, 24, 22, 17, 10, 20, 15, 9}
		for _, val := range values {
			sc.Update(val)
		}
		sc.Update(15)
		a.Equalf(false, sc.Long(), "Indicator should not recommend to buy")
		a.Equalf(false, sc.Short(), "Indicator should not recommend to sell")
	}
}
