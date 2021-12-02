package indicator

import "sync"

type MACDEvaluator struct {
	mu sync.RWMutex // mutex to protect macd evaluator

	emaS *EMAEvaluator // EMA with short period
	emaL *EMAEvaluator // EMA with long period
	emaA *EMAEvaluator // EMA with short period for smoothing MACD

	macd   float64 // MACD value
	signal float64 // MACD signal line value
}

func NewMACDEvaluator(shortPeriod, longPeriod, averagePeriod int, alpha AlphaFunc) *MACDEvaluator {
	return &MACDEvaluator{
		emaS: NewEMAEvaluator(shortPeriod, alpha),
		emaL: NewEMAEvaluator(longPeriod, alpha),
		emaA: NewEMAEvaluator(averagePeriod, alpha),
	}
}

func (m *MACDEvaluator) UpdateMACD(p float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.emaS.UpdateEMA(p)
	m.emaL.UpdateEMA(p)
	m.macd = m.emaS.GetEMA() - m.emaL.GetEMA()

	m.emaA.UpdateEMA(m.macd)
	m.signal = m.emaA.GetEMA()
}

func (m *MACDEvaluator) GetMACD() (macd float64, signal float64) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.macd, m.signal
}

type MACDStrategy struct {
	mu sync.RWMutex // mutex to protect strategy

	macd       *MACDEvaluator
	prevMACD   float64
	prevSignal float64
}

func NewMACDStrategy(macd *MACDEvaluator) Strategy {
	return &MACDStrategy{
		macd: macd,
	}
}

func (m *MACDStrategy) Update(p float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.prevMACD, m.prevSignal = m.macd.GetMACD()
	m.macd.UpdateMACD(p)
}

func (m *MACDStrategy) Long() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	curMACD, curSignal := m.macd.GetMACD()
	return m.prevMACD < m.prevSignal && curMACD > curSignal
}

func (m *MACDStrategy) Short() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	curMACD, curSignal := m.macd.GetMACD()
	return m.prevMACD > m.prevSignal && curMACD < curSignal
}
