package main

import (
	"github.com/dataprism/dataprism-commons/api"
	//"github.com/dataprism/dataprism-logics/logics"
	consul "github.com/hashicorp/consul/api"
	//nomad "github.com/hashicorp/nomad/api"
	"github.com/sirupsen/logrus"
	"flag"
	//"github.com/dataprism/dataprism-logics/evals"
	//"github.com/dataprism/dataprism-lib/nodes"
	"strconv"
	"github.com/dataprism/dataprism-sync/connectors"
	consul2 "github.com/dataprism/dataprism-commons/consul"
	"github.com/dataprism/dataprism-sync/links"
)

func main() {
	//var jobsDir = flag.String("d", "/tmp", "the directory where job information will be stored")
	var port = flag.Int("p", 6400, "the port of the dataprism logics rest api")

	API := api.CreateAPI("0.0.0.0:" + strconv.Itoa(*port))

	client, err := consul.NewClient(consul.DefaultConfig())
	if err != nil {
		logrus.Error(err)
	}

	storage := consul2.NewStorage(client)

	//nomadClient, err := nomad.NewClient(nomad.DefaultConfig())
	//if err != nil {
	//	logrus.Error(err)
	//}

	// -- Connectors
	connectorManager := connectors.NewManager(storage)
	connectorRouter := connectors.NewRouter(connectorManager)
	API.RegisterGet("/v1/connectors", connectorRouter.ListConnectors)
	API.RegisterGet("/v1/connectors/{id}", connectorRouter.GetConnector)
	API.RegisterPost("/v1/connectors", connectorRouter.SetConnector)
	API.RegisterDelete("/v1/connectors/{id}", connectorRouter.RemoveConnector)

	// -- Links
	linkManager := links.NewManager(storage)
	linkRouter := links.NewRouter(linkManager)
	API.RegisterGet("/v1/links", linkRouter.ListLinks)
	API.RegisterGet("/v1/links/{id}", linkRouter.GetLink)
	API.RegisterPost("/v1/links", linkRouter.SetLink)
	API.RegisterDelete("/v1/links/{id}", linkRouter.RemoveLink)


	err = API.Start()
	if err != nil {
		logrus.Error(err)
	}
}