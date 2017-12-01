package scheduler

import (
	"context"
	"github.com/dataprism/dataprism-commons/schedule"
	"github.com/dataprism/dataprism-sync/links"
	"github.com/dataprism/dataprism-sync/connectors"
)

type ExecutionManager struct {
	linkManager *links.LinkManager
	connectorManager *connectors.ConnectorManager
	scheduler *schedule.Scheduler
}

func NewManager(linkManager *links.LinkManager, connectorManager *connectors.ConnectorManager, scheduler *schedule.Scheduler) *ExecutionManager {
	return &ExecutionManager{linkManager, connectorManager, scheduler}
}

func (m *ExecutionManager) Deploy(ctx context.Context, id string) (*schedule.ScheduleResponse, error) {
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
	job := NewSyncJob(link, connector)

	// -- schedule the job
	return m.scheduler.Schedule(job)
}