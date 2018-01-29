package sync

import (
	"github.com/hashicorp/nomad/api"
	"github.com/dataprism/dataprism-commons/utils"
	"strings"
	"strconv"
	"errors"
	"encoding/json"
	"github.com/hashicorp/nomad/helper"
	"github.com/dataprism/dataprism-commons/config"
	"github.com/dataprism/dataprism-sync-runtime/plugins"
)

type SyncJob struct {
	link *Link
	input *plugins.InputType
	output *plugins.OutputType
	cluster *config.KafkaCluster
}

func DefaultSyncJobResources() *api.Resources {
	return &api.Resources{
		CPU:      helper.IntToPtr(1500),
		MemoryMB: helper.IntToPtr(512),
		IOPS:     helper.IntToPtr(0),
	}
}

func NewSyncJob(link *Link, input *plugins.InputType, output *plugins.OutputType, cluster *config.KafkaCluster) SyncJob {
	return SyncJob{link, input, output, cluster}
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

	task := api.NewTask(nomadJobId, "docker")

	task.Config = make(map[string]interface{})
	task.Config["image"] = "dataprism/dataprism-sync-runtime:latest"
	task.Env = vars
	task.Meta = make(map[string]string)
	task.Meta["link"] = string(strLink)
	task.Resources = DefaultSyncJobResources()

	taskGroup := api.NewTaskGroup("syncs", 1)
	taskGroup.Tasks = []*api.Task{task}

	nomadJob := api.NewServiceJob(nomadJobId, strings.ToTitle(s.link.Name), "global", 1)
	nomadJob.Datacenters = []string{ "aws" }
	nomadJob.TaskGroups = []*api.TaskGroup{taskGroup}

	return nomadJob, nil
}

func (s *SyncJob) generateEnvVars(cluster *config.KafkaCluster) (map[string]string, error) {

	result := make(map[string]string)

	// -- input
	for k, v := range s.link.InputSettings { result["input." + s.link.InputTypeId + "." + k] = v }
	for k, v := range s.link.OutputSettings { result["output." + s.link.InputTypeId + "." + k] = v }

	if s.link.OutputTypeId == "kafka-output" {
		result["output.type"] = "kafka-output"
		result["output.kafka-output.bootstrap_servers"] = strings.Join(cluster.Servers, ",")
		result["output.kafka-output.min.messages"] = strconv.Itoa(cluster.KafkaBufferMinMsg)
		result["output.kafka-output.buffering.max.ms"] = strconv.Itoa(cluster.KafkaBufferMaxMs)
		result["output.kafka-output.data_topic"] = s.link.Topic
	}

	if s.link.InputTypeId == "kafka-input" {
		result["input.type"] = "kafka-input"
		result["input.kafka-input.bootstrap_servers"] = strings.Join(cluster.Servers, ",")
		result["input.kafka-input.topic"] = s.link.Topic
		result["input.kafka-input.group_id"] = s.link.Id

	}

	return result, nil;
}