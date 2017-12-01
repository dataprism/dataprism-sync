package main

import (
	"github.com/dataprism/dataprism-commons/api"
	consul "github.com/hashicorp/consul/api"
	nomad "github.com/hashicorp/nomad/api"
	"github.com/sirupsen/logrus"
	"flag"
	"strconv"
	"github.com/dataprism/dataprism-sync/connectors"
	consul2 "github.com/dataprism/dataprism-commons/consul"
	"github.com/dataprism/dataprism-sync/links"
	"github.com/dataprism/dataprism-commons/schedule"
	"github.com/dataprism/dataprism-sync/scheduler"
)

func main() {
	var jobsDir = flag.String("d", "/tmp", "the directory where job information will be stored")
	var port = flag.Int("p", 6400, "the port of the dataprism sync rest api")

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
	connectorManager := connectors.NewManager(storage)
	connectorRouter := connectors.NewRouter(connectorManager)
	API.RegisterGet("/v1/connectors", connectorRouter.ListConnectors)
	API.RegisterGet("/v1/connectors/{id}", connectorRouter.GetConnector)
	API.RegisterPost("/v1/connectors", connectorRouter.SetConnector)
	API.RegisterDelete("/v1/connectors/{id}", connectorRouter.RemoveConnector)

	// -- Links
	linkManager := links.NewManager(storage, s)
	linkRouter := links.NewRouter(linkManager)
	API.RegisterGet("/v1/links", linkRouter.ListLinks)
	API.RegisterGet("/v1/links/{id}", linkRouter.GetLink)
	API.RegisterPost("/v1/links", linkRouter.SetLink)
	API.RegisterDelete("/v1/links/{id}", linkRouter.RemoveLink)

	// -- Executions
	executionManager := scheduler.NewManager(linkManager, connectorManager, s)
	executionRouter := scheduler.NewRouter(executionManager)
	API.RegisterPost("/v1/links/{id}/run", executionRouter.Deploy)

	err = API.Start()
	if err != nil {
		logrus.Error(err)
	}
}