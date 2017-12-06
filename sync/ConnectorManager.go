package sync

import (
	"context"
	"encoding/json"
	"github.com/dataprism/dataprism-commons/core"
)

type ConnectorManager struct {
	platform *core.Platform
}

func NewConnectorManager(platform *core.Platform) *ConnectorManager {
	return &ConnectorManager{platform}
}

func (m *ConnectorManager) ListConnectors(ctx context.Context) ([]*Connector, error) {
	var result []*Connector

	pairs, err := m.platform.KV.List(ctx, "connectors/")
	if err != nil { return nil, err }

	for _, p := range pairs {
		var entity Connector
		if err = json.Unmarshal(p.Value, &entity); err != nil { return nil, err }
		result = append(result, &entity)
	}

	return result, err
}

func (m *ConnectorManager) GetConnector(ctx context.Context, id string) (*Connector, error) {
	data, err := m.platform.KV.Get(ctx, "connectors/" + id)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, nil
	}

	var res Connector
	err = json.Unmarshal(data.Value, &res)
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

	err = m.platform.KV.Set(ctx, "connectors/" + connector.Id, data)
	if err != nil { return nil, err }
	return connector, nil
}

func (m *ConnectorManager) RemoveConnector(ctx context.Context, id string) (error) {
	return m.platform.KV.Remove(ctx, "connectors/" + id)
}