package console

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"buscharge/internal/audit"
	"buscharge/internal/bus"
	"buscharge/internal/depot"
	"buscharge/internal/grid"
	"buscharge/internal/pile"
	"buscharge/internal/plan"
	"buscharge/internal/power"
	"buscharge/internal/quota"
	"buscharge/internal/soc"
	"buscharge/internal/store"
)

type Config struct {
	Addr    string
	DataDir string
	WebDir  string
}

type Server struct {
	addr   string
	webDir string
	store  *store.Store
	soc    *soc.Service
	grid   *grid.Capacity
	quota  *quota.Accumulator
	limits *quota.Limits
	buses  *bus.Service
	power  *power.Service
	piles  *pile.Service
	plans  *plan.Service
	depot  *depot.Service
	audit  *audit.Service
}

func NewServer(cfg Config) (*Server, error) {
	st, err := store.New(cfg.DataDir)
	if err != nil {
		return nil, err
	}
	s := &Server{
		addr:   cfg.Addr,
		webDir: cfg.WebDir,
		store:  st,
		soc:    soc.NewService(),
		grid:   grid.NewCapacity(80000),
		quota:  quota.NewAccumulator(),
		limits: quota.NewLimits(),
		buses:  bus.NewService(),
	}
	s.power = power.NewService(s.grid, s.quota, st)
	s.piles = pile.NewService(s.soc, s.buses, s.power, st)
	s.plans = plan.NewService(s.piles, s.soc, st)
	s.depot = depot.NewService("depot-main", s.piles, s.buses, s.grid)
	s.audit = audit.NewService(st)
	return s, nil
}

func (s *Server) Start() error {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req, "/buses", http.StatusFound)
	})
	r.Get("/buses", s.servePage("buses"))
	r.Get("/piles", s.servePage("piles"))
	r.Get("/power", s.servePage("power"))
	r.Get("/alarms", s.servePage("alarms"))

	r.Route("/api", func(api chi.Router) {
		api.Get("/health", s.handleHealth)
		api.Get("/status", s.handleStatus)
		api.Get("/dashboard", s.handleDashboard)
		api.Get("/report", s.handleReport)
		api.Get("/insights", s.handleInsights)
		api.Get("/snapshot", s.handleSnapshot)
		api.Post("/charge/cycle", s.handleChargeCycle)
		api.Post("/charge/terminate", s.handleTerminationCycle)
		api.Get("/namespaces", s.handleNamespaces)
		api.Post("/zones", s.handleAddZone)
		api.Post("/zones/piles", s.handleZonePile)
		api.Get("/zones/{zone}/piles", s.handleZonePiles)

		api.Get("/buses", s.handleBuses)
		api.Get("/buses/stats", s.handleBusStats)
		api.Post("/buses", s.handleAddBus)
		api.Post("/buses/{id}/authorize", s.handleAuthorizeBus)
		api.Post("/buses/route", s.handleBusRoute)

		api.Get("/piles", s.handlePiles)
		api.Get("/piles/stored/{id}", s.handleStoredPile)
		api.Get("/piles/sessions", s.handlePileSessions)
		api.Post("/piles", s.handleRegisterPile)
		api.Post("/piles/assign", s.handleAssignPile)
		api.Post("/piles/plug", s.handlePlugPile)
		api.Post("/piles/allocate", s.handleAllocatePile)
		api.Post("/piles/renumber", s.handleRenumberPile)
		api.Post("/piles/{id}/start", s.handleStartPile)
		api.Post("/piles/{id}/stop", s.handleStopPile)
		api.Post("/piles/{id}/charge", s.handleChargePile)
		api.Post("/piles/{id}/terminate", s.handleTerminatePile)
		api.Post("/piles/{id}/session", s.handleBeginSession)
		api.Delete("/piles/{id}/session", s.handleEndSession)
		api.Post("/piles/duration", s.handlePileDuration)

		api.Get("/plans/{id}", s.handlePlanLoad)
		api.Post("/plans/read", s.handlePlanRead)
		api.Post("/plans/read-context", s.handlePlanReadContext)
		api.Post("/plans/update", s.handlePlanUpdate)
		api.Post("/plans/assign", s.handlePlanAssign)
		api.Post("/plans/schedule", s.handlePlanSchedule)
		api.Post("/plans/window", s.handlePlanWindow)
		api.Get("/plans/conflicts", s.handlePlanConflicts)

		api.Get("/power", s.handlePower)
		api.Post("/power/limit", s.handlePowerLimit)
		api.Post("/power/coordinate", s.handlePowerCoordinate)
		api.Post("/power/meter", s.handlePowerMeter)
		api.Get("/power/meter/{id}", s.handleStoredMeter)
		api.Post("/power/shed", s.handlePowerShed)
		api.Post("/power/balance", s.handlePowerBalance)
		api.Post("/power/demand", s.handleDemand)

		api.Get("/grid", s.handleGrid)
		api.Get("/grid/feeders", s.handleGridFeeders)
		api.Post("/grid/expand", s.handleGridExpand)
		api.Post("/grid/set", s.handleGridSet)
		api.Post("/grid/feeders", s.handleGridAddFeeder)

		api.Get("/quota", s.handleQuota)
		api.Post("/quota/set", s.handleQuotaSet)
		api.Get("/quota/limits/{id}", s.handleQuotaLimitGet)
		api.Post("/quota/limits", s.handleQuotaLimitSet)
		api.Post("/quota/allow", s.handleQuotaAllow)

		api.Get("/soc/{id}", s.handleSocGet)
		api.Post("/soc", s.handleSocSet)
		api.Post("/soc/{id}/refresh", s.handleSocRefresh)
		api.Post("/soc/{id}/energy", s.handleSocEnergy)

		api.Get("/audit", s.handleAuditList)
		api.Get("/audit/filter", s.handleAuditFilter)
		api.Get("/audit/pile", s.handleAuditByPile)
		api.Get("/audit/bus", s.handleAuditByBus)
		api.Get("/audit/depot", s.handleAuditByDepot)
		api.Get("/audit/{id}", s.handleAuditGet)
		api.Post("/audit", s.handleAuditRecord)
		api.Delete("/audit/{id}", s.handleAuditDelete)
	})

	return http.ListenAndServe(s.addr, r)
}
