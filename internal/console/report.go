package console

import (
	"net/http"

	"buscharge/internal/audit"
	"buscharge/internal/bus"
	"buscharge/internal/depot"
	"buscharge/internal/pile"
	"buscharge/internal/plan"
)

type DepotReport struct {
	Stats       depot.Stats            `json:"stats"`
	Fleet       bus.FleetStats         `json:"fleet"`
	Roster      []bus.RosterEntry      `json:"roster"`
	Piles       []pile.Pile            `json:"piles"`
	Sessions    []pile.Session         `json:"sessions"`
	ChargeOrder []plan.ChargeCandidate `json:"charge_order"`
	Utilization float64                `json:"utilization"`
	PeakLoad    int                    `json:"peak_load"`
	AverageSoc  int                    `json:"average_soc"`
	Audit       audit.Summary          `json:"audit"`
	Projected   int                    `json:"projected"`
}

func (s *Server) buildReport(loads []int) (DepotReport, error) {
	report := DepotReport{
		Stats:       s.depot.Stats(),
		Fleet:       s.buses.FleetStats(),
		Roster:      s.buses.Roster(),
		Piles:       s.piles.List(),
		Sessions:    s.piles.ActiveSessions(),
		ChargeOrder: s.plans.ChargeOrder(s.busIDs()),
		Utilization: s.power.Utilization(),
		PeakLoad:    s.power.PeakLoad(loads),
		AverageSoc:  s.soc.AverageSoc(s.busIDs()),
	}
	if ids := s.busIDs(); len(ids) > 0 {
		report.Projected = s.soc.ProjectSoc(ids[0], 60000)
	}
	summary, err := s.audit.Summary()
	if err != nil {
		return DepotReport{}, err
	}
	report.Audit = summary
	return report, nil
}

func (s *Server) busIDs() []string {
	ids := make([]string, 0)
	for _, v := range s.buses.List() {
		ids = append(ids, v.ID)
	}
	return ids
}

func (s *Server) handleReport(w http.ResponseWriter, req *http.Request) {
	report, err := s.buildReport([]int{90000, 30000, 120000})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, report)
}

func (s *Server) handleAuditFilter(w http.ResponseWriter, req *http.Request) {
	event := req.URL.Query().Get("event")
	records, err := s.audit.Filter(event)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, records)
}
