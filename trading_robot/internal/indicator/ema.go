package indicator

type EMAEvaluator struct {
	counter int     // value counter
	ema     float64 // EMA value
	alpha   float64 // alpha coefficient
}

type AlphaFunc func(period int) float64

func NewEMAEvaluator(period int, alpha AlphaFunc) *EMAEvaluator {
	return &EMAEvaluator{
		ema:   0,
		alpha: alpha(period),
	}
}

// EMA is exponential moving average. The EMA for a series P may be calculated recursively:
// EMA(1) = P(1)                                    t = 0
// EMA(t) = alpha * P(t) + (1 - alpha) * EMA(t-1)   t > 0
// The coefficient alpha represents the degree of weighting decrease, a constant smoothing factor between 0 and 1.
// A higher alpha discounts older observations faster. P(t) is the value at a time period t.
func EMA(pEMA float64, p float64, alpha float64) float64 {
	return alpha*p + (1-alpha)*pEMA
}

func (e *EMAEvaluator) UpdateEMA(p float64) {
	e.counter++
	if e.counter == 1 {
		e.ema = p
	}
	e.ema = EMA(e.ema, p, e.alpha)
}

func (e *EMAEvaluator) GetEMA() float64 {
	return e.ema
}

type EMAStrategy struct {
	ema      *EMAEvaluator
	curPrice float64
}

func NewEMAStrategy(ema *EMAEvaluator) Strategy {
	return &EMAStrategy{
		ema: ema,
	}
}

func (e *EMAStrategy) Update(p float64) {
	e.ema.UpdateEMA(p)
	e.curPrice = p
}

func (e *EMAStrategy) Long() bool {
	return e.curPrice > e.ema.GetEMA()
}

func (e *EMAStrategy) Short() bool {
	return e.curPrice < e.ema.GetEMA()
}
