package main

import (
	"github.com/dataprism/dataprism-commons/api"
	consul "github.com/hashicorp/consul/api"
	nomad "github.com/hashicorp/nomad/api"
	"github.com/sirupsen/logrus"
	"flag"
	"strconv"
	consul2 "github.com/dataprism/dataprism-commons/consul"
	"github.com/dataprism/dataprism-commons/schedule"
	"strings"
	"github.com/dataprism/dataprism-sync/core"
)

func main() {
	var jobsDir = flag.String("d", "/tmp", "the directory where job information will be stored")
	var port = flag.Int("p", 6400, "the port of the dataprism sync rest api")

	var kafkaServers = flag.String("kafka-servers", "localhost:9092", "the kafka cluster nodes")
	var kafkaBufferMaxMs = flag.Int("kafka-buffer-max-ms", 1000, "the max amount of time to buffer events before sending them to kafka")
	var kafkaBufferMinMsg = flag.Int("kafka-buffer-min-msg", 1000, "the min amount of messages to buffer events before sending them to kafka")

	cluster := &core.KafkaCluster{strings.Split(*kafkaServers, ","), *kafkaBufferMaxMs, *kafkaBufferMinMsg}

	API := api.CreateAPI("0.0.0.0:" + strconv.Itoa(*port))

	client, err := consul.NewClient(consul.DefaultConfig())
	if err != nil {
		logrus.Error(err)
	}

	storage := consul2.NewStorage(client)

	nomadClient, err := nomad.NewClient(nomad.DefaultConfig())
	if err != nil {
		logrus.Error(err)
	}

	s := schedule.NewScheduler(nomadClient, *jobsDir)

	// -- Connectors
	connectorManager := core.NewConnectorManager(storage)
	connectorRouter := core.NewConnectorRouter(connectorManager)
	API.RegisterGet("/v1/connectors", connectorRouter.ListConnectors)
	API.RegisterGet("/v1/connectors/{id}", connectorRouter.GetConnector)
	API.RegisterPost("/v1/connectors", connectorRouter.SetConnector)
	API.RegisterDelete("/v1/connectors/{id}", connectorRouter.RemoveConnector)

	// -- Links
	linkManager := core.NewLinkManager(storage, s)
	linkRouter := core.NewLinkRouter(linkManager)
	API.RegisterGet("/v1/links", linkRouter.ListLinks)
	API.RegisterGet("/v1/links/{id}", linkRouter.GetLink)
	API.RegisterPost("/v1/links", linkRouter.SetLink)
	API.RegisterDelete("/v1/links/{id}", linkRouter.RemoveLink)

	// -- Executions
	executionManager := core.NewExecutionManager(linkManager, connectorManager, s, cluster)
	executionRouter := core.NewExecutionRouter(executionManager)
	API.RegisterPost("/v1/links/{id}/run", executionRouter.Deploy)

	err = API.Start()
	if err != nil {
		logrus.Error(err)
	}
}