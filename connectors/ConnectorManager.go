package connectors

import (
	"context"
	"encoding/json"
	"github.com/dataprism/dataprism-commons/consul"
)

type ConnectorManager struct {
	consulStorage *consul.ConsulStorage
}

func NewManager(consulStorage *consul.ConsulStorage) *ConnectorManager {
	return &ConnectorManager{consulStorage: consulStorage}
}

func (m *ConnectorManager) ListConnectors(ctx context.Context) ([]string, error) {
	return m.consulStorage.List(ctx, "connectors/")
}

func (m *ConnectorManager) GetConnector(ctx context.Context, id string) (*Connector, error) {
	data, err := m.consulStorage.Get(ctx, "connectors/" + id)
	if err != nil {
		return nil, err
	}

	var res Connector
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (m *ConnectorManager) SetConnector(ctx context.Context, connector *Connector) (*Connector, error) {
	data, err := json.Marshal(connector)
	if err != nil {
		return nil, err
	}

	err = m.consulStorage.Set(ctx, "connectors/" + connector.Id, data)
	if err != nil { return nil, err }
	return connector, nil
}

func (m *ConnectorManager) RemoveConnector(ctx context.Context, id string) (error) {
	return m.consulStorage.Remove(ctx, "connectors/" + id)
}