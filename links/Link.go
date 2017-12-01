package links

import "github.com/dataprism/dataprism-sync/cluster"

type Link struct {
	Id string						`json:"id"`
	Name string						`json:"name"`
	Direction string				`json:"direction"`
	ConnectorId string				`json:"connector_id"`
	WorkersPerProcess int			`json:"workers_per_process"`
	Processes int					`json:"processes"`
	Settings map[string]string		`json:"settings"`
	KafkaCluster *cluster.KafkaCluster `json:"cluster"`
	Topic string  					`json:"topic"`
}