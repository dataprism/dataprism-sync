package sync

type Link struct {
	Id string							`json:"id"`
	Name string							`json:"name"`
	Direction string					`json:"direction"`
	InputTypeId string					`json:"input_type_id"`
	InputSettings map[string]string		`json:"input_settings"`
	OutputTypeId string					`json:"output_type_id"`
	OutputSettings map[string]string	`json:"input_settings"`
	WorkersPerProcess int				`json:"workers_per_process"`
	Processes int						`json:"processes"`
	Topic string  						`json:"topic"`
}