package verifycase

import (
	"testing"

	"buscharge/internal/bus"
	"buscharge/internal/grid"
	"buscharge/internal/pile"
	"buscharge/internal/power"
	"buscharge/internal/quota"
	"buscharge/internal/soc"
	"buscharge/internal/store"
)

func TestBcChargeTerminationOrder(t *testing.T) {
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
	s.SetCurrent("b1", 100)
	if err := ps.Start("p1", 20); err != nil {
		t.Fatal(err)
	}
	verdict := ps.TerminateCharge("b1", "p1")
	if !verdict.Full {
		t.Fatal("full verdict should be raised at the charge boundary")
	}
	if !verdict.TrickleStopped {
		t.Fatal("trickle stop must apply before the full verdict")
	}
}
