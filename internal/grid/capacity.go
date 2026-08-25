package grid

import "sync"

type Capacity struct {
	mu      sync.RWMutex
	watts   int
	feeders map[string]int
}

func NewCapacity(watts int) *Capacity {
	return &Capacity{watts: watts, feeders: make(map[string]int)}
}

func (c *Capacity) Current() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.watts
}

func (c *Capacity) Expand(additional int) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.watts += additional
	return c.watts
}

func (c *Capacity) Set(watts int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.watts = watts
}
