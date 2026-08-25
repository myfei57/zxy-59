package grid

type Feeder struct {
	Name  string `json:"name"`
	Limit int    `json:"limit"`
}

func (c *Capacity) AddFeeder(name string, limit int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.feeders[name] = limit
}

func (c *Capacity) Feeders() []Feeder {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make([]Feeder, 0, len(c.feeders))
	for name, limit := range c.feeders {
		out = append(out, Feeder{Name: name, Limit: limit})
	}
	return out
}

func (c *Capacity) Feeder(name string) (Feeder, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	limit, ok := c.feeders[name]
	return Feeder{Name: name, Limit: limit}, ok
}
