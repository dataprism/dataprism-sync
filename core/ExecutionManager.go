package core

import (
	"context"
	"github.com/dataprism/dataprism-commons/schedule"
)

type ExecutionManager struct {
	linkManager *LinkManager
	connectorManager *ConnectorManager
	scheduler *schedule.Scheduler
	cluster *KafkaCluster
}

func NewExecutionManager(linkManager *LinkManager, connectorManager *ConnectorManager, scheduler *schedule.Scheduler, cluster *KafkaCluster) *ExecutionManager {
	return &ExecutionManager{linkManager, connectorManager, scheduler, cluster}
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
	job := &SyncJob{link, connector, m.cluster}

	// -- schedule the job
	return m.scheduler.Schedule(job)
}

func (m *ExecutionManager) Undeploy(ctx context.Context, id string) (*schedule.UnscheduleResponse, error) {
	// -- schedule the job
	return m.scheduler.Unschedule("sync", id)
}