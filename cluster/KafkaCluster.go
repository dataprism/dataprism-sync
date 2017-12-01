package cluster

type KafkaCluster struct {
	Servers []string				`json:"servers"`
	KafkaBufferMaxMs int			`json:"buffer_max_ms"`
	KafkaBufferMinMsg int			`json:"buffer_min_msg"`
}