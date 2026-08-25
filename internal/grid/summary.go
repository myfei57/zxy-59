package grid

func (c *Capacity) FeederTotal() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	total := 0
	for _, limit := range c.feeders {
		total += limit
	}
	return total
}

func (c *Capacity) FeederHeadroom() int {
	return c.Current() - c.FeederTotal()
}

func (c *Capacity) FeederCount() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.feeders)
}
