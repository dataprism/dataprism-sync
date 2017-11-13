package links

type Link struct {
	Id string
	Name string
	ConnectorId string
	WorkersPerProcess int
	Processes int
	Settings map[string]string
}
