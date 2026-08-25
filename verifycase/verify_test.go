package verifycase

import (
	"os"
	"testing"

	"buscharge/internal/bus"
	"buscharge/internal/grid"
	"buscharge/internal/pile"
	"buscharge/internal/plan"
	"buscharge/internal/power"
	"buscharge/internal/quota"
	"buscharge/internal/soc"
	"buscharge/internal/store"
)

func TestBcPlanDurableFirst(t *testing.T) {
	dir := t.TempDir()
	st, err := store.New(dir)
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
	if err := ps.Register("p2", "slow"); err != nil {
		t.Fatal(err)
	}
	plans := plan.NewService(ps, s, st)
	blockPath := st.Path(st.PlanKey("b3"))
	if err := os.MkdirAll(blockPath, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := plans.Update("b3", "p1", 1); err == nil {
		t.Fatal("plan update must fail when the plan cannot be persisted")
	}
	if owner := ps.Owner("p1"); owner == "b3" {
		t.Fatal("pile must not switch when the plan is not durable")
	}
}
