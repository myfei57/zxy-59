package console

import (
	"net/http"

	"buscharge/internal/power"
)

func (s *Server) handleSnapshot(w http.ResponseWriter, req *http.Request) {
	auditSummary, _ := s.audit.Summary()
	writeJSON(w, http.StatusOK, map[string]any{
		"stats":    s.depot.Stats(),
		"zones":    s.depot.ListZones(),
		"fleet":    s.buses.FleetStats(),
		"roster":   s.buses.Roster(),
		"piles":    s.piles.List(),
		"sessions": s.piles.ActiveSessions(),
		"capacity": s.depot.GridCapacity(),
		"quota":    s.power.QuotaTotal(),
		"feeders":  s.grid.Feeders(),
		"audit":    auditSummary,
	})
}

func (s *Server) handleAuditByPile(w http.ResponseWriter, req *http.Request) {
	pileID := req.URL.Query().Get("pile")
	records, err := s.audit.ByPile(pileID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, records)
}

func (s *Server) handleAuditByBus(w http.ResponseWriter, req *http.Request) {
	busID := req.URL.Query().Get("bus")
	records, err := s.audit.ByBus(busID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, records)
}

func (s *Server) handleAuditByDepot(w http.ResponseWriter, req *http.Request) {
	depot := req.URL.Query().Get("depot")
	records, err := s.audit.ByDepot(depot)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, records)
}

func (s *Server) handleDemand(w http.ResponseWriter, req *http.Request) {
	var body struct {
		Requests []int `json:"requests"`
	}
	if err := readJSON(req, &body); err != nil {
		writeError(w, http.StatusBadRequest, "requests is required")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"total":   power.TotalRequested(body.Requests),
		"average": power.AverageRequested(body.Requests),
	})
}
