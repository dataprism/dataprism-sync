package core

import (
	"github.com/hashicorp/nomad/api"
	"github.com/dataprism/dataprism-commons/utils"
	"strings"
	"strconv"
	"errors"
	"encoding/json"
	"github.com/hashicorp/nomad/helper"
)

type SyncJob struct {
	link *Link
	connector *Connector
	cluster *KafkaCluster
}

func DefaultSyncJobResources() *api.Resources {
	return &api.Resources{
		CPU:      helper.IntToPtr(1500),
		MemoryMB: helper.IntToPtr(512),
		IOPS:     helper.IntToPtr(0),
	}
}

func NewSyncJob(link *Link, connector *Connector, cluster *KafkaCluster) SyncJob {
	return SyncJob{link, connector, cluster}
}

func (s *SyncJob) ToJob() (*api.Job, error) {
	nomadJobId := utils.ToNomadJobId("sync", s.link.Id)

	vars, err := s.generateEnvVars(s.cluster)
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

	task.Config = make(map[string]interface{})
	task.Config["image"] = "dataprism/dataprism-sync-runtime:latest"
	task.Env = vars
	task.Meta = make(map[string]string)
	task.Meta["link"] = string(strLink)
	task.Meta["connector"] = string(strConnector)
	task.Resources = DefaultSyncJobResources()

	taskGroup := api.NewTaskGroup("syncs", 1)
	taskGroup.Tasks = []*api.Task{task}

	nomadJob := api.NewServiceJob(nomadJobId, strings.ToTitle(s.link.Name), "global", 1)
	nomadJob.Datacenters = []string{ "aws" }
	nomadJob.TaskGroups = []*api.TaskGroup{taskGroup}

	return nomadJob, nil
}

func (s *SyncJob) generateEnvVars(cluster *KafkaCluster) (map[string]string, error) {

	result := make(map[string]string)

	if s.link.Direction == "IN" {
		if !s.connector.IsInput {
			return nil, errors.New("the connector " + s.connector.Name + "(" + s.connector.Id + ") is not an input connector")
		}

		result["input.type"] = s.connector.Type
		for k, v := range s.link.Settings {
			result["input." + s.connector.Type + "." + k] = v
		}

		result["output.type"] = "kafka"
		result["output.kafka.bootstrap_servers"] = strings.Join(cluster.Servers, ",")
		result["output.kafka.min.messages"] = strconv.Itoa(cluster.KafkaBufferMinMsg)
		result["output.kafka.buffering.max.ms"] = strconv.Itoa(cluster.KafkaBufferMaxMs)
		result["output.kafka.data_topic"] = s.link.Topic

	} else if s.link.Direction == "OUT" {
		if !s.connector.IsOutput {
			return nil, errors.New("the connector " + s.connector.Name + "(" + s.connector.Id + ") is not an output connector")
		}

		result["input.type"] = "kafka"
		result["input.kafka.bootstrap_servers"] = strings.Join(cluster.Servers, ",")
		result["input.kafka.topic"] = s.link.Topic
		result["input.kafka.group_id"] = s.link.Id

		result["output.type"] = s.connector.Type
		for k, v := range s.link.Settings {
			result["output." + s.connector.Type + "." + k] = v
		}

	} else {
		return nil, errors.New("invalid link direction " + s.link.Direction)
	}

	return result, nil;
}