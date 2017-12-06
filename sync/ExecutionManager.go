package sync

import (
	"context"
	"github.com/dataprism/dataprism-commons/execute"
	"github.com/dataprism/dataprism-commons/core"
)

type ExecutionManager struct {
	platform *core.Platform
	linkManager *LinkManager
	connectorManager *ConnectorManager
}

func NewExecutionManager(platform *core.Platform, linkManager *LinkManager, connectorManager *ConnectorManager) *ExecutionManager {
	return &ExecutionManager{platform,linkManager, connectorManager}
}

func (m *ExecutionManager) Deploy(ctx context.Context, id string) (*execute.ScheduleResponse, error) {
	// -- get the link
	link, err := m.linkManager.GetLink(ctx, id)
	if err != nil {
		return nil, err
	}

	// -- get the connector for the link
	connector, err := m.connectorManager.GetConnector(ctx, link.ConnectorId)
	if err != nil {
		return nil, err
	}

	// -- create the job for the link
	job := &SyncJob{link, connector, m.platform.Settings.KafkaCluster}

	// -- schedule the job
	return m.platform.Scheduler.Schedule(job)
}

func (m *ExecutionManager) Undeploy(ctx context.Context, id string) (*execute.UnscheduleResponse, error) {
	// -- schedule the job
	return m.platform.Scheduler.Unschedule("sync", id)
}