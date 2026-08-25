package console

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"buscharge/internal/ns"
)

func (s *Server) handleHealth(w http.ResponseWriter, req *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (s *Server) handleStatus(w http.ResponseWriter, req *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"data_dir": s.store.Root(),
		"piles":    len(s.piles.List()),
		"buses":    len(s.buses.List()),
	})
}

func (s *Server) handleNamespaces(w http.ResponseWriter, req *http.Request) {
	zones := s.depot.ListZones()
	out := make([]map[string]string, 0, len(zones))
	for _, n := range zones {
		out = append(out, map[string]string{
			"depot":     n.Depot,
			"zone":      n.Zone,
			"key":       n.Key(),
			"depot_key": ns.DepotKey(n.Depot),
			"zone_key":  ns.ZoneKey(n.Depot, n.Zone),
		})
	}
	writeJSON(w, http.StatusOK, out)
}

func (s *Server) handleAddZone(w http.ResponseWriter, req *http.Request) {
	var body struct {
		Zone string `json:"zone"`
	}
	if err := readJSON(req, &body); err != nil || body.Zone == "" {
		writeError(w, http.StatusBadRequest, "zone is required")
		return
	}
	n := s.depot.AddZone(body.Zone)
	writeJSON(w, http.StatusOK, n)
}

func (s *Server) handleBuses(w http.ResponseWriter, req *http.Request) {
	vehicles := s.buses.List()
	type item struct {
		ID        string `json:"id"`
		Route     string `json:"route"`
		Pile      string `json:"pile"`
		Departure string `json:"departure"`
	}
	out := make([]item, 0, len(vehicles))
	for _, v := range vehicles {
		out = append(out, item{
			ID:        v.ID,
			Route:     v.Route,
			Pile:      s.buses.Mapping(v.ID),
			Departure: s.buses.Departure(v.ID),
		})
	}
	writeJSON(w, http.StatusOK, out)
}

func (s *Server) handleAddBus(w http.ResponseWriter, req *http.Request) {
	var body struct {
		ID      string `json:"id"`
		Route   string `json:"route"`
		Depart  string `json:"departure"`
		Allowed bool   `json:"authorized"`
	}
	if err := readJSON(req, &body); err != nil || body.ID == "" {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}
	s.buses.AddVehicle(body.ID, body.Route)
	s.buses.SetDeparture(body.ID, body.Depart)
	if body.Allowed {
		s.buses.Authorize(body.ID)
	}
	vehicle, _ := s.buses.Vehicle(body.ID)
	writeJSON(w, http.StatusOK, vehicle)
}

func (s *Server) handleAuthorizeBus(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	s.buses.Authorize(id)
	writeJSON(w, http.StatusOK, map[string]any{"id": id, "authorized": true})
}

func (s *Server) handlePiles(w http.ResponseWriter, req *http.Request) {
	writeJSON(w, http.StatusOK, s.piles.List())
}

func (s *Server) handleRegisterPile(w http.ResponseWriter, req *http.Request) {
	var body struct {
		Number string `json:"number"`
		Kind   string `json:"kind"`
	}
	if err := readJSON(req, &body); err != nil || body.Number == "" {
		writeError(w, http.StatusBadRequest, "number is required")
		return
	}
	if err := s.piles.Register(body.Number, body.Kind); err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"number": body.Number})
}

func (s *Server) handleAssignPile(w http.ResponseWriter, req *http.Request) {
	var body struct {
		BusID string `json:"bus_id"`
	}
	if err := readJSON(req, &body); err != nil {
		writeError(w, http.StatusBadRequest, "bus_id is required")
		return
	}
	number, err := s.piles.Assign(body.BusID)
	if err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"number": number})
}

func (s *Server) handlePlugPile(w http.ResponseWriter, req *http.Request) {
	var body struct {
		BusID  string `json:"bus_id"`
		Number string `json:"number"`
	}
	if err := readJSON(req, &body); err != nil {
		writeError(w, http.StatusBadRequest, "bus_id and number are required")
		return
	}
	if err := s.piles.Plug(body.BusID, body.Number); err != nil {
		writeError(w, http.StatusForbidden, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"engaged": s.piles.Engaged(body.Number)})
}

func (s *Server) handleAllocatePile(w http.ResponseWriter, req *http.Request) {
	var body struct {
		BusID  string `json:"bus_id"`
		Number string `json:"number"`
	}
	if err := readJSON(req, &body); err != nil {
		writeError(w, http.StatusBadRequest, "bus_id and number are required")
		return
	}
	ok := s.piles.Allocate(body.Number, body.BusID)
	writeJSON(w, http.StatusOK, map[string]any{"allocated": ok, "owner": s.piles.Owner(body.Number)})
}

func (s *Server) handleRenumberPile(w http.ResponseWriter, req *http.Request) {
	var body struct {
		BusID     string `json:"bus_id"`
		OldNumber string `json:"old_number"`
		NewNumber string `json:"new_number"`
	}
	if err := readJSON(req, &body); err != nil {
		writeError(w, http.StatusBadRequest, "bus_id, old_number and new_number are required")
		return
	}
	if err := s.depot.RenumberPile(body.BusID, body.OldNumber, body.NewNumber); err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"mapping": s.buses.Mapping(body.BusID)})
}

func (s *Server) handleStartPile(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	var body struct {
		Requested int `json:"requested"`
	}
	if err := readJSON(req, &body); err != nil {
		writeError(w, http.StatusBadRequest, "requested is required")
		return
	}
	if err := s.piles.Start(id, body.Requested); err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	snap, _ := s.piles.Snapshot(id)
	writeJSON(w, http.StatusOK, snap)
}

func (s *Server) handleStopPile(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	if err := s.piles.Stop(id); err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	snap, _ := s.piles.Snapshot(id)
	writeJSON(w, http.StatusOK, snap)
}

func (s *Server) handleChargePile(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	var body struct {
		Requested int `json:"requested"`
	}
	if err := readJSON(req, &body); err != nil {
		writeError(w, http.StatusBadRequest, "requested is required")
		return
	}
	allowed := s.piles.BeginCharge(id, body.Requested)
	writeJSON(w, http.StatusOK, map[string]any{"pile": id, "allowed": allowed})
}

func (s *Server) handleTerminatePile(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	var body struct {
		BusID string `json:"bus_id"`
	}
	if err := readJSON(req, &body); err != nil {
		writeError(w, http.StatusBadRequest, "bus_id is required")
		return
	}
	writeJSON(w, http.StatusOK, s.piles.TerminateCharge(body.BusID, id))
}

func (s *Server) handlePlanRead(w http.ResponseWriter, req *http.Request) {
	var body struct {
		BusID string `json:"bus_id"`
	}
	if err := readJSON(req, &body); err != nil {
		writeError(w, http.StatusBadRequest, "bus_id is required")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"percent": s.plans.Read(body.BusID)})
}

func (s *Server) handlePlanReadContext(w http.ResponseWriter, req *http.Request) {
	var body struct {
		BusID string `json:"bus_id"`
	}
	if err := readJSON(req, &body); err != nil {
		writeError(w, http.StatusBadRequest, "bus_id is required")
		return
	}
	percent, err := s.plans.ReadWithContext(body.BusID, req.Context())
	if err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"percent": percent})
}

func (s *Server) handlePlanUpdate(w http.ResponseWriter, req *http.Request) {
	var body struct {
		BusID    string `json:"bus_id"`
		PileID   string `json:"pile_id"`
		Sequence int    `json:"sequence"`
	}
	if err := readJSON(req, &body); err != nil {
		writeError(w, http.StatusBadRequest, "bus_id and pile_id are required")
		return
	}
	if err := s.plans.Update(body.BusID, body.PileID, body.Sequence); err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"updated": true})
}

func (s *Server) handlePlanLoad(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	p, err := s.plans.Load(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (s *Server) handlePlanAssign(w http.ResponseWriter, req *http.Request) {
	var body struct {
		BusID  string `json:"bus_id"`
		Number string `json:"number"`
	}
	if err := readJSON(req, &body); err != nil {
		writeError(w, http.StatusBadRequest, "bus_id and number are required")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"allocated": s.plans.Assign(body.BusID, body.Number)})
}

func (s *Server) handlePlanSchedule(w http.ResponseWriter, req *http.Request) {
	var body struct {
		BusIDs []string `json:"bus_ids"`
	}
	if err := readJSON(req, &body); err != nil {
		writeError(w, http.StatusBadRequest, "bus_ids is required")
		return
	}
	writeJSON(w, http.StatusOK, s.plans.Schedule(body.BusIDs))
}

func (s *Server) handlePower(w http.ResponseWriter, req *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"capacity": s.depot.GridCapacity(),
		"quota":    s.power.QuotaTotal(),
		"reserved": s.power.ReservedTotal(),
	})
}

func (s *Server) handlePowerLimit(w http.ResponseWriter, req *http.Request) {
	var body struct {
		Requested int `json:"requested"`
	}
	if err := readJSON(req, &body); err != nil {
		writeError(w, http.StatusBadRequest, "requested is required")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"allowed": s.power.Limit(body.Requested)})
}

func (s *Server) handlePowerCoordinate(w http.ResponseWriter, req *http.Request) {
	var body struct {
		Requests []int `json:"requests"`
	}
	if err := readJSON(req, &body); err != nil {
		writeError(w, http.StatusBadRequest, "requests is required")
		return
	}
	writeJSON(w, http.StatusOK, s.power.Coordinate(body.Requests))
}

func (s *Server) handlePowerMeter(w http.ResponseWriter, req *http.Request) {
	var body struct {
		PileID  string  `json:"pile_id"`
		Reading float64 `json:"reading"`
	}
	if err := readJSON(req, &body); err != nil {
		writeError(w, http.StatusBadRequest, "pile_id and reading are required")
		return
	}
	writeJSON(w, http.StatusOK, s.power.Meter(body.PileID, body.Reading))
}

func (s *Server) handleGrid(w http.ResponseWriter, req *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"capacity": s.depot.GridCapacity()})
}

func (s *Server) handleGridExpand(w http.ResponseWriter, req *http.Request) {
	var body struct {
		Watts int `json:"watts"`
	}
	if err := readJSON(req, &body); err != nil {
		writeError(w, http.StatusBadRequest, "watts is required")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"capacity": s.depot.ExpandGrid(body.Watts)})
}

func (s *Server) handleGridSet(w http.ResponseWriter, req *http.Request) {
	var body struct {
		Watts int `json:"watts"`
	}
	if err := readJSON(req, &body); err != nil {
		writeError(w, http.StatusBadRequest, "watts is required")
		return
	}
	s.grid.Set(body.Watts)
	writeJSON(w, http.StatusOK, map[string]any{"capacity": s.grid.Current()})
}

func (s *Server) handleQuota(w http.ResponseWriter, req *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"total": s.power.QuotaTotal()})
}

func (s *Server) handleQuotaSet(w http.ResponseWriter, req *http.Request) {
	var body struct {
		Total float64 `json:"total"`
	}
	if err := readJSON(req, &body); err != nil {
		writeError(w, http.StatusBadRequest, "total is required")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"total": s.quota.Set(body.Total)})
}

func (s *Server) handleSocGet(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	writeJSON(w, http.StatusOK, map[string]any{
		"current":  s.soc.Current(id),
		"estimate": s.soc.Estimate(id),
	})
}

func (s *Server) handleSocSet(w http.ResponseWriter, req *http.Request) {
	var body struct {
		BusID   string `json:"bus_id"`
		Percent int    `json:"percent"`
	}
	if err := readJSON(req, &body); err != nil || body.BusID == "" {
		writeError(w, http.StatusBadRequest, "bus_id and percent are required")
		return
	}
	s.soc.SetCurrent(body.BusID, body.Percent)
	writeJSON(w, http.StatusOK, map[string]any{"current": s.soc.Current(body.BusID)})
}

func (s *Server) handleSocRefresh(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	writeJSON(w, http.StatusOK, map[string]any{"estimate": s.soc.Refresh(id)})
}

func (s *Server) handleAuditGet(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	if !s.store.Has("audit/" + id) {
		writeError(w, http.StatusNotFound, "audit record not found")
		return
	}
	rec, err := s.audit.Get(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, rec)
}

func (s *Server) handleAuditRecord(w http.ResponseWriter, req *http.Request) {
	var body struct {
		Depot  string  `json:"depot"`
		Event  string  `json:"event"`
		BusID  string  `json:"bus_id"`
		PileID string  `json:"pile_id"`
		Amount float64 `json:"amount"`
	}
	if err := readJSON(req, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid payload")
		return
	}
	rec, err := s.audit.Record(body.Depot, body.Event, body.BusID, body.PileID, body.Amount)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, rec)
}

func (s *Server) handleAuditDelete(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	if err := s.store.Delete("audit/" + id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"deleted": true})
}
