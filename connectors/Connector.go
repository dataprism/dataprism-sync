package connectors

type Connector struct {
	Id string
	Name string
	Properties []Property
}

type Property struct {
	Key string			`json:"key"`
	Description string	`json:"description"`
	Priority string		`json:"priority"`
}
