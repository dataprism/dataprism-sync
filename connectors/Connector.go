package connectors

type Connector struct {
	Id string				`json:"id"`
	Name string				`json:"name"`
	Properties []Property	`json:"properties"`
}

type Property struct {
	Key string			`json:"key"`
	Description string	`json:"description"`
	Priority string		`json:"priority"`
}
