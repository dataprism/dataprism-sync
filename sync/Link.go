package sync

type Link struct {
	Id string						`json:"id"`
	Name string						`json:"name"`
	Direction string				`json:"direction"`
	ConnectorId string				`json:"connector_id"`
	WorkersPerProcess int			`json:"workers_per_process"`
	Processes int					`json:"processes"`
	Settings map[string]string		`json:"settings"`
	Topic string  					`json:"topic"`
}