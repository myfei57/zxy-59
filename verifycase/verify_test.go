package verifycase

import (
	"testing"

	"buscharge/internal/bus"
	"buscharge/internal/depot"
	"buscharge/internal/grid"
	"buscharge/internal/pile"
	"buscharge/internal/power"
	"buscharge/internal/quota"
	"buscharge/internal/soc"
	"buscharge/internal/store"
)

func TestBcZonePileAlias(t *testing.T) {
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
	d := depot.NewService("d", ps, b, g)
	d.AddZone("z1")
	d.AssignPileToZone("z1", "p1")
	d.AssignPileToZone("z1", "p2")
	piles := d.ZonePiles("z1")
	piles[0] = "p9"
	after := d.ZonePiles("z1")
	if after[0] != "p1" {
		t.Fatalf("depot state must not be corrupted by callers, got %s", after[0])
	}
}
