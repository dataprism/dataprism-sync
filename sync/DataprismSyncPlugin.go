package sync

import (
	"github.com/dataprism/dataprism-commons/core"
	"github.com/dataprism/dataprism-commons/api"
)

type DataprismPlugin struct {}

func (d *DataprismPlugin) Id() string { return "sync" }

func (d *DataprismPlugin) CreateRoutes(platform *core.Platform, API *api.Rest) {
	// -- Connectors
	connectorManager := NewConnectorManager(platform)
	connectorRouter := NewConnectorRouter(connectorManager)
	API.RegisterGet("/v1/connectors", connectorRouter.ListConnectors)
	API.RegisterGet("/v1/connectors/{id}", connectorRouter.GetConnector)
	API.RegisterPost("/v1/connectors", connectorRouter.SetConnector)
	API.RegisterDelete("/v1/connectors/{id}", connectorRouter.RemoveConnector)

	// -- Links
	linkManager := NewLinkManager(platform)
	linkRouter := NewLinkRouter(linkManager)
	API.RegisterGet("/v1/links", linkRouter.ListLinks)
	API.RegisterGet("/v1/links/{id}", linkRouter.GetLink)
	API.RegisterPost("/v1/links", linkRouter.SetLink)
	API.RegisterDelete("/v1/links/{id}", linkRouter.RemoveLink)

	// -- Executions
	executionManager := NewExecutionManager(platform, linkManager, connectorManager)
	executionRouter := NewExecutionRouter(executionManager)
	API.RegisterPost("/v1/links/{id}/run", executionRouter.Deploy)
	API.RegisterDelete("/v1/links/{id}/run", executionRouter.Undeploy)
}