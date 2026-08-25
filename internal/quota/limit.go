package quota

import "sync"

type Limits struct {
	mu    sync.RWMutex
	daily map[string]float64
	used  map[string]float64
}

func NewLimits() *Limits {
	return &Limits{
		daily: make(map[string]float64),
		used:  make(map[string]float64),
	}
}

func (l *Limits) SetDaily(pileID string, limit float64) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.daily[pileID] = limit
}

func (l *Limits) Daily(pileID string) float64 {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.daily[pileID]
}

func (l *Limits) Used(pileID string) float64 {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.used[pileID]
}

func (l *Limits) Allow(pileID string, amount float64) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	limit := l.daily[pileID]
	if limit <= 0 {
		return true
	}
	if l.used[pileID]+amount > limit {
		return false
	}
	l.used[pileID] += amount
	return true
}
