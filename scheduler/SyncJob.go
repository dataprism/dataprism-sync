package scheduler

import (
	"github.com/hashicorp/nomad/api"
	"github.com/dataprism/dataprism-sync/links"
	"github.com/dataprism/dataprism-sync/connectors"
	"github.com/dataprism/dataprism-commons/utils"
	"strings"
	"strconv"
	"errors"
	"encoding/json"
)

type SyncJob struct {
	link *links.Link
	connector *connectors.Connector
}

func NewSyncJob(link *links.Link, connector *connectors.Connector) *SyncJob {
	return &SyncJob{link, connector}
}

func (s *SyncJob) ToJob() (*api.Job, error) {
	nomadJobId := utils.ToNomadJobId("sync", s.link.Id)

	vars, err := s.generateEnvVars()
	if err != nil {
		return nil, err
	}

	strLink, err := json.Marshal(s.link)
	if err != nil {
		return nil, errors.New("the link could not be serialized to the job metadata")
	}

	strConnector, err := json.Marshal(s.connector)
	if err != nil {
		return nil, errors.New("the connector could not be serialized to the job metadata")
	}

	task := api.NewTask(nomadJobId, "docker")

	//task.Config = make(map[string]interface{})
	task.Config["image"] = "vrtoeni/to-kafka:latest"
	task.Env = vars
	task.Meta["link"] = string(strLink)
	task.Meta["connector"] = string(strConnector)

	taskGroup := api.NewTaskGroup("syncs", 1)
	taskGroup.Tasks = []*api.Task{task}

	nomadJob := api.NewServiceJob(nomadJobId, strings.ToTitle(s.link.Name), "global", 1)
	nomadJob.Datacenters = []string{ "dc1" }
	nomadJob.TaskGroups = []*api.TaskGroup{taskGroup}

	return nomadJob, nil
}

func (s *SyncJob) generateEnvVars() (map[string]string, error) {

	result := make(map[string]string)

	if s.link.Direction == "IN" {
		if !s.connector.IsInput {
			return nil, errors.New("the connector " + s.connector.Name + "(" + s.connector.Id + ") is not an input connector")
		}

		result["input.type"] = s.connector.Type
		for k, v := range s.link.Settings {
			result["input." + k] = v
		}

		result["output.type"] = "kafka"
		result["output.kafka.bootstrap_servers"] = strings.Join(s.link.KafkaCluster.Servers, ",")
		result["output.kafka.min.messages"] = strconv.Itoa(s.link.KafkaCluster.KafkaBufferMinMsg)
		result["output.kafka.buffering.max.ms"] = strconv.Itoa(s.link.KafkaCluster.KafkaBufferMaxMs)
		result["output.kafka.data_topic"] = s.link.Topic

	} else if s.link.Direction == "OUT" {
		if !s.connector.IsOutput {
			return nil, errors.New("the connector " + s.connector.Name + "(" + s.connector.Id + ") is not an output connector")
		}

		result["input.type"] = "kafka"
		result["input.kafka.bootstrap_servers"] = strings.Join(s.link.KafkaCluster.Servers, ",")
		result["input.kafka.topic"] = s.link.Topic
		result["input.kafka.group_id"] = s.link.Id

		result["output.type"] = s.connector.Type
		for k, v := range s.link.Settings {
			result["output." + k] = v
		}

	} else {
		return nil, errors.New("invalid link direction " + s.link.Direction)
	}

	return result, nil;
}