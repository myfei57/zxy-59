package quota

import "sync"

type Accumulator struct {
	mu    sync.RWMutex
	total float64
}

func NewAccumulator() *Accumulator {
	return &Accumulator{}
}

func (a *Accumulator) Total() float64 {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.total
}

func (a *Accumulator) Set(value float64) float64 {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.total = value
	return a.total
}

func (a *Accumulator) Add(amount float64) float64 {
	a.total += amount
	return a.total
}
