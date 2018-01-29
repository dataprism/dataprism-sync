package sync

import (
	"github.com/dataprism/dataprism-commons/core"
	"github.com/dataprism/dataprism-commons/api"
	"github.com/dataprism/dataprism-sync-runtime/plugins"
	"github.com/dataprism/dataprism-sync-runtime/plugins/elasticsearch"
	"github.com/dataprism/dataprism-sync-runtime/plugins/kafka"
	"github.com/dataprism/dataprism-sync-runtime/plugins/rest"
)

type DataprismPlugin struct {}

func (d *DataprismPlugin) Id() string { return "sync" }

func (d *DataprismPlugin) CreateRoutes(platform *core.Platform, API *api.Rest) {
	// -- Connectors
	registry := plugins.NewSyncPluginRegistry()
	registry.Add(kafka.NewKafkaSyncPlugin())
	registry.Add(elasticsearch.NewElasticsearchSyncPlugin())
	registry.Add(rest.NewRestSyncPlugin())

	pluginRouter := NewPluginRouter(registry)

	API.RegisterGet("/v1/sync/plugins", pluginRouter.List)
	API.RegisterGet("/v1/sync/plugins/inputs", pluginRouter.ListInputs)
	API.RegisterGet("/v1/sync/plugins/outputs", pluginRouter.ListOutputs)

	// -- we currently disable adding, editing and removing connectors
	//API.RegisterPost("/v1/connectors", connectorRouter.SetConnector)
	//API.RegisterDelete("/v1/connectors/{id}", connectorRouter.RemoveConnector)

	// -- Links
	linkManager := NewLinkManager(platform)
	linkRouter := NewLinkRouter(linkManager)
	API.RegisterGet("/v1/sync/links", linkRouter.ListLinks)
	API.RegisterGet("/v1/sync/links/{id}", linkRouter.GetLink)
	API.RegisterPost("/v1/sync/links", linkRouter.SetLink)
	API.RegisterDelete("/v1/sync/links/{id}", linkRouter.RemoveLink)

	// -- Executions
	executionManager := NewExecutionManager(platform, linkManager, registry)
	executionRouter := NewExecutionRouter(executionManager)
	API.RegisterPost("/v1/sync/links/{id}/run", executionRouter.Deploy)
	API.RegisterDelete("/v1/sync/links/{id}/run", executionRouter.Undeploy)
}