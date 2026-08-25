package verifycase

import (
	"sync"
	"testing"

	"buscharge/internal/grid"
	"buscharge/internal/power"
	"buscharge/internal/quota"
	"buscharge/internal/store"
)

func TestBcConcurrentPowerMeter(t *testing.T) {
	st, err := store.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	g := grid.NewCapacity(80000)
	q := quota.NewAccumulator()
	pw := power.NewService(g, q, st)
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			pw.Meter("p1", 1.0)
		}()
	}
	wg.Wait()
	if total := pw.QuotaTotal(); total != 100.0 {
		t.Fatalf("accumulated charge must include every meter read, got %v", total)
	}
}
