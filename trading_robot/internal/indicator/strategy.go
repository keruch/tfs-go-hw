package indicator

type Strategy interface {
	Update(p float64)
	Long() bool
	Short() bool
}

type StrategiesComposition []Strategy

func NewStrategiesComposition(strategies ...Strategy) Strategy {
	sc := make(StrategiesComposition, 0)
	sc = append(sc, strategies...)
	return sc
}

func (sc StrategiesComposition) Update(p float64) {
	for _, strategy := range sc {
		strategy.Update(p)
	}
}

func (sc StrategiesComposition) Long() bool {
	var (
		long  = true
		short = false
	)
	for _, strategy := range sc {
		long = long && strategy.Long()
		short = short || strategy.Short()
	}

	return long && !short
}

func (sc StrategiesComposition) Short() bool {
	var (
		long  = false
		short = true
	)
	for _, strategy := range sc {
		short = short && strategy.Short()
		long = long || strategy.Long()
	}

	return short && !long
}
