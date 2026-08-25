package ns

func DepotKey(depot string) string {
	return "depot/" + depot
}

func ZoneKey(depot, zone string) string {
	return "depot/" + depot + "/zone/" + zone
}
