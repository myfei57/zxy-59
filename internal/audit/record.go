package audit

import "time"

type Record struct {
	ID     string    `json:"id"`
	Depot  string    `json:"depot"`
	Event  string    `json:"event"`
	BusID  string    `json:"bus_id"`
	PileID string    `json:"pile_id"`
	Amount float64   `json:"amount"`
	At     time.Time `json:"at"`
}
