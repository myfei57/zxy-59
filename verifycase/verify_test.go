package verifycase

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"

	"buscharge/internal/bus"
	"buscharge/internal/grid"
	"buscharge/internal/pile"
	"buscharge/internal/power"
	"buscharge/internal/quota"
	"buscharge/internal/soc"
	"buscharge/internal/store"
)

func TestBcConcurrentPileAlloc(t *testing.T) {
	st, err := store.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	s := soc.NewService()
	b := bus.NewService()
	g := grid.NewCapacity(80000)
	q := quota.NewAccumulator()
	pw := power.NewService(g, q, st)
	ps := pile.NewService(s, b, pw, st)
	if err := ps.Register("p1", "fast"); err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	var successes int32
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			if ps.Allocate("p1", fmt.Sprintf("bus-%d", n)) {
				atomic.AddInt32(&successes, 1)
			}
		}(i)
	}
	wg.Wait()
	if atomic.LoadInt32(&successes) != 1 {
		t.Fatalf("a pile must be allocated to exactly one bus, got %d", successes)
	}
}
