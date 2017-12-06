package sync

import (
	"context"
	"encoding/json"
	"github.com/dataprism/dataprism-commons/core"
)

type LinkManager struct {
	platform *core.Platform
}

func NewLinkManager(platform *core.Platform) *LinkManager {
	return &LinkManager{platform}
}

func (m *LinkManager) ListLinks(ctx context.Context) ([]*Link, error) {
	var result []*Link

	pairs, err := m.platform.KV.List(ctx, "links/")
	if err != nil { return nil, err }

	for _, p := range pairs {
		var entity Link
		if err = json.Unmarshal(p.Value, &entity); err != nil { return nil, err }
		result = append(result, &entity)
	}

	return result, err
}

func (m *LinkManager) GetLink(ctx context.Context, id string) (*Link, error) {
	data, err := m.platform.KV.Get(ctx, "links/" + id)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, nil
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

	err = m.platform.KV.Set(ctx, "links/" + link.Id, data)
	if err != nil { return nil, err }
	return link, nil
}

func (m *LinkManager) RemoveLink(ctx context.Context, id string) (error) {
	return m.platform.KV.Remove(ctx, "links/" + id)
}