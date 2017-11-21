package links

type Link struct {
	Id string						`json:"id"`
	Name string						`json:"name"`
	ConnectorId string				`json:"connector_id"`
	WorkersPerProcess int			`json:"workers_per_process"`
	Processes int					`json:"processes"`
	Settings map[string]string		`json:"settings"`
}
