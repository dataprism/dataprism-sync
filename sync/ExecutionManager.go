package sync

import (
	"context"
	"github.com/dataprism/dataprism-commons/execute"
	"github.com/dataprism/dataprism-commons/core"
	"github.com/dataprism/dataprism-sync-runtime/plugins"
	"github.com/pkg/errors"
)

type ExecutionManager struct {
	platform *core.Platform
	linkManager *LinkManager
	pluginRegistry *plugins.PluginRegistry
}

func NewExecutionManager(platform *core.Platform, linkManager *LinkManager, pluginRegistry *plugins.PluginRegistry) *ExecutionManager {
	return &ExecutionManager{platform,linkManager, pluginRegistry}
}

func (m *ExecutionManager) Deploy(ctx context.Context, id string) (*execute.ScheduleResponse, error) {
	// -- get the link
	link, err := m.linkManager.GetLink(ctx, id)
	if err != nil {
		return nil, err
	}

	// -- get the plugin for the link
	inputType, ok := m.pluginRegistry.GetInputType(link.InputTypeId)
	if !ok {
		return nil, errors.New("The input type with id " + link.InputTypeId + " could not be found")
	}

	outputType, ok := m.pluginRegistry.GetOutputType(link.OutputTypeId)
	if !ok {
		return nil, errors.New("The output type with id " + link.OutputTypeId + " could not be found")
	}

	// -- create the job for the link
	job := &SyncJob{link, inputType, outputType, m.platform.Settings.KafkaCluster}

	// -- schedule the job
	return m.platform.Scheduler.Schedule(job)
}

func (m *ExecutionManager) Undeploy(ctx context.Context, id string) (*execute.UnscheduleResponse, error) {
	// -- schedule the job
	return m.platform.Scheduler.Unschedule("sync", id)
}