package sync

type Connector struct {
	Id string				`json:"id"`
	Type string				`json:"type"`
	Name string				`json:"name"`
	IsOutput bool			`json:"is_output"`
	IsInput bool			`json:"is_input"`
	Properties []Property	`json:"properties"`
}

type Property struct {
	Key string			`json:"key"`
	Description string	`json:"description"`
	Priority string		`json:"priority"`
}
