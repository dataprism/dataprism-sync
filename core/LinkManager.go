package core

import (
	"context"
	"encoding/json"
	"github.com/dataprism/dataprism-commons/consul"
	"github.com/dataprism/dataprism-commons/schedule"
)

type LinkManager struct {
	storage *consul.ConsulStorage
	scheduler *schedule.Scheduler
}

func NewLinkManager(consulStorage *consul.ConsulStorage, scheduler *schedule.Scheduler) *LinkManager {
	return &LinkManager{storage: consulStorage, scheduler:scheduler}
}

func (m *LinkManager) ListLinks(ctx context.Context) ([]*Link, error) {
	var result []*Link

	pairs, err := m.storage.List(ctx, "links/")
	if err != nil { return nil, err }

	for _, p := range pairs {
		var entity Link
		if err = json.Unmarshal(p.Value, &entity); err != nil { return nil, err }
		result = append(result, &entity)
	}

	return result, err
}

func (m *LinkManager) GetLink(ctx context.Context, id string) (*Link, error) {
	data, err := m.storage.Get(ctx, "links/" + id)
	if err != nil {
		return nil, err
	}

	var res Link
	err = json.Unmarshal(data.Value, &res)
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

	err = m.storage.Set(ctx, "links/" + link.Id, data)
	if err != nil { return nil, err }
	return link, nil
}

func (m *LinkManager) RemoveLink(ctx context.Context, id string) (error) {
	return m.storage.Remove(ctx, "links/" + id)
}