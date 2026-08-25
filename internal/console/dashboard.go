package console

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"buscharge/internal/pile"
	"buscharge/internal/power"
)

func (s *Server) handleDashboard(w http.ResponseWriter, req *http.Request) {
	auditSummary, _ := s.audit.Summary()
	writeJSON(w, http.StatusOK, map[string]any{
		"fleet":           s.buses.FleetStats(),
		"zones":           s.depot.ZoneCount(),
		"routes":          s.buses.RouteCount(),
		"active_sessions": s.piles.ActiveCount(),
		"plan_windows":    s.plans.WindowCount(),
		"piles":           s.piles.List(),
		"capacity":        s.depot.GridCapacity(),
		"quota":           s.power.QuotaTotal(),
		"audit":           auditSummary,
		"feeders":         s.grid.Feeders(),
		"rates":           s.pileRates(),
		"stored": map[string]int{
			"plans":  s.store.Count("plan"),
			"piles":  s.store.Count("pile"),
			"meters": s.store.Count("meter"),
			"audit":  s.store.Count("audit"),
		},
	})
}

func (s *Server) pileRates() map[string]int {
	out := make(map[string]int)
	for _, p := range s.piles.List() {
		out[p.Number] = s.piles.Rate(p.Number)
	}
	return out
}

func (s *Server) handleBusStats(w http.ResponseWriter, req *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"stats":  s.buses.FleetStats(),
		"routes": s.buses.RouteCount(),
	})
}

func (s *Server) handleBusRoute(w http.ResponseWriter, req *http.Request) {
	var body struct {
		Route  string   `json:"route"`
		BusIDs []string `json:"bus_ids"`
	}
	if err := readJSON(req, &body); err != nil || body.Route == "" {
		writeError(w, http.StatusBadRequest, "route is required")
		return
	}
	s.buses.SetRoute(body.Route, body.BusIDs)
	writeJSON(w, http.StatusOK, map[string]any{"route": body.Route, "buses": s.buses.Route(body.Route)})
}

func (s *Server) handleZonePile(w http.ResponseWriter, req *http.Request) {
	var body struct {
		Zone   string `json:"zone"`
		Number string `json:"number"`
	}
	if err := readJSON(req, &body); err != nil || body.Zone == "" || body.Number == "" {
		writeError(w, http.StatusBadRequest, "zone and number are required")
		return
	}
	s.depot.AssignPileToZone(body.Zone, body.Number)
	writeJSON(w, http.StatusOK, map[string]any{"piles": s.depot.ZonePiles(body.Zone)})
}

func (s *Server) handleZonePiles(w http.ResponseWriter, req *http.Request) {
	zone := chi.URLParam(req, "zone")
	writeJSON(w, http.StatusOK, map[string]any{"zone": zone, "piles": s.depot.ZonePiles(zone)})
}

func (s *Server) handlePileSessions(w http.ResponseWriter, req *http.Request) {
	writeJSON(w, http.StatusOK, s.piles.ActiveSessions())
}

func (s *Server) handleBeginSession(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	var body struct {
		BusID string `json:"bus_id"`
	}
	if err := readJSON(req, &body); err != nil {
		writeError(w, http.StatusBadRequest, "bus_id is required")
		return
	}
	if err := s.piles.BeginSession(id, body.BusID); err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	session, _ := s.piles.Session(id)
	writeJSON(w, http.StatusOK, session)
}

func (s *Server) handleEndSession(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	if err := s.piles.EndSession(id); err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	session, _ := s.piles.Session(id)
	writeJSON(w, http.StatusOK, session)
}

func (s *Server) handlePileDuration(w http.ResponseWriter, req *http.Request) {
	var body struct {
		Kind      string `json:"kind"`
		WattHours int    `json:"watt_hours"`
	}
	if err := readJSON(req, &body); err != nil {
		writeError(w, http.StatusBadRequest, "kind and watt_hours are required")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"power":    pile.RatedPower(body.Kind),
		"duration": pile.EstimateDuration(body.Kind, body.WattHours),
	})
}

func (s *Server) handlePlanWindow(w http.ResponseWriter, req *http.Request) {
	var body struct {
		BusID string `json:"bus_id"`
		Start string `json:"start"`
		End   string `json:"end"`
	}
	if err := readJSON(req, &body); err != nil {
		writeError(w, http.StatusBadRequest, "bus_id, start and end are required")
		return
	}
	start, err := time.Parse(time.RFC3339, body.Start)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid start")
		return
	}
	end, err := time.Parse(time.RFC3339, body.End)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid end")
		return
	}
	s.plans.SetWindow(body.BusID, start, end)
	window, _ := s.plans.Window(body.BusID)
	writeJSON(w, http.StatusOK, window)
}

func (s *Server) handlePlanConflicts(w http.ResponseWriter, req *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"conflicts": s.plans.Conflicts()})
}

func (s *Server) handlePowerShed(w http.ResponseWriter, req *http.Request) {
	var body struct {
		Requests []int `json:"requests"`
	}
	if err := readJSON(req, &body); err != nil {
		writeError(w, http.StatusBadRequest, "requests is required")
		return
	}
	writeJSON(w, http.StatusOK, s.power.Shed(body.Requests))
}

func (s *Server) handlePowerBalance(w http.ResponseWriter, req *http.Request) {
	var body struct {
		Requests []int `json:"requests"`
	}
	if err := readJSON(req, &body); err != nil {
		writeError(w, http.StatusBadRequest, "requests is required")
		return
	}
	phases := power.BalancePhases(body.Requests)
	writeJSON(w, http.StatusOK, map[string]any{
		"phases":    phases,
		"imbalance": power.PhaseImbalance(phases[0], phases[1], phases[2]),
	})
}

func (s *Server) handleGridFeeders(w http.ResponseWriter, req *http.Request) {
	writeJSON(w, http.StatusOK, s.grid.Feeders())
}

func (s *Server) handleGridAddFeeder(w http.ResponseWriter, req *http.Request) {
	var body struct {
		Name  string `json:"name"`
		Limit int    `json:"limit"`
	}
	if err := readJSON(req, &body); err != nil || body.Name == "" {
		writeError(w, http.StatusBadRequest, "name and limit are required")
		return
	}
	s.grid.AddFeeder(body.Name, body.Limit)
	feeder, _ := s.grid.Feeder(body.Name)
	writeJSON(w, http.StatusOK, feeder)
}

func (s *Server) handleQuotaLimitGet(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	writeJSON(w, http.StatusOK, map[string]any{
		"daily": s.limits.Daily(id),
		"used":  s.limits.Used(id),
	})
}

func (s *Server) handleQuotaLimitSet(w http.ResponseWriter, req *http.Request) {
	var body struct {
		PileID string  `json:"pile_id"`
		Daily  float64 `json:"daily"`
	}
	if err := readJSON(req, &body); err != nil || body.PileID == "" {
		writeError(w, http.StatusBadRequest, "pile_id and daily are required")
		return
	}
	s.limits.SetDaily(body.PileID, body.Daily)
	writeJSON(w, http.StatusOK, map[string]any{"pile_id": body.PileID, "daily": s.limits.Daily(body.PileID)})
}

func (s *Server) handleQuotaAllow(w http.ResponseWriter, req *http.Request) {
	var body struct {
		PileID string  `json:"pile_id"`
		Amount float64 `json:"amount"`
	}
	if err := readJSON(req, &body); err != nil {
		writeError(w, http.StatusBadRequest, "pile_id and amount are required")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"allowed": s.limits.Allow(body.PileID, body.Amount)})
}

func (s *Server) handleSocEnergy(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	var body struct {
		Capacity  int `json:"capacity"`
		Energy    int `json:"energy"`
		Charge    int `json:"charge"`
		Discharge int `json:"discharge"`
	}
	if err := readJSON(req, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid payload")
		return
	}
	if body.Capacity > 0 {
		s.soc.SetCapacityWh(id, body.Capacity)
	}
	if body.Energy > 0 {
		s.soc.SetEnergy(id, body.Energy)
	}
	if body.Charge != 0 {
		s.soc.Charge(id, body.Charge)
	}
	if body.Discharge != 0 {
		s.soc.Discharge(id, body.Discharge)
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"capacity": s.soc.CapacityWh(id),
		"energy":   s.soc.Energy(id),
		"current":  s.soc.Current(id),
	})
}

func (s *Server) handleAuditList(w http.ResponseWriter, req *http.Request) {
	records, err := s.audit.List()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	summary, _ := s.audit.Summary()
	writeJSON(w, http.StatusOK, map[string]any{"records": records, "summary": summary})
}

func (s *Server) handleStoredPile(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	var stored pile.Pile
	if err := s.store.LoadPile(id, &stored); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, stored)
}

func (s *Server) handleStoredMeter(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	var stored power.MeterReading
	if err := s.store.LoadMeter(id, &stored); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, stored)
}
