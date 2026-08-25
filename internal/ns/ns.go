package ns

type Namespace struct {
	Depot string `json:"depot"`
	Zone  string `json:"zone"`
}

func New(depot, zone string) Namespace {
	return Namespace{Depot: depot, Zone: zone}
}

func (n Namespace) Key() string {
	return n.Depot + "/" + n.Zone
}
