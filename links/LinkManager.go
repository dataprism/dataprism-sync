package links

import (
	"context"
	"encoding/json"
	"github.com/dataprism/dataprism-lib/consul"
)

type LinkManager struct {
	consulStorage *consul.ConsulStorage
}

func NewManager(consulStorage *consul.ConsulStorage) *LinkManager {
	return &LinkManager{consulStorage: consulStorage}
}

func (m *LinkManager) ListLinks(ctx context.Context) ([]string, error) {
	return m.consulStorage.List(ctx, "links/")
}

func (m *LinkManager) GetLink(ctx context.Context, id string) (*Link, error) {
	data, err := m.consulStorage.Get(ctx, "links/" + id)
	if err != nil {
		return nil, err
	}

	var res Link
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (m *LinkManager) SetLink(ctx context.Context, link *Link) (*Link, error) {
	data, err := json.Marshal(link)
	if err != nil {
		return nil, err
	}

	err = m.consulStorage.Set(ctx, "links/" + link.Id, data)
	if err != nil { return nil, err }
	return link, nil
}

func (m *LinkManager) RemoveLink(ctx context.Context, id string) (error) {
	return m.consulStorage.Remove(ctx, "links/" + id)
}