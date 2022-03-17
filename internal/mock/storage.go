package mock

import "github.com/DuarteMRAlves/maestro/internal"

type AssetStorage struct {
	Assets map[internal.AssetName]internal.Asset
}

func (m AssetStorage) Save(asset internal.Asset) error {
	m.Assets[asset.Name()] = asset
	return nil
}

func (m AssetStorage) Load(name internal.AssetName) (
	internal.Asset,
	error,
) {
	asset, exists := m.Assets[name]
	if !exists {
		err := &internal.NotFound{Type: "asset", Ident: name.Unwrap()}
		return internal.Asset{}, err
	}
	return asset, nil
}

type OrchestrationStorage struct {
	Orchs map[internal.OrchestrationName]internal.Orchestration
}

func (m OrchestrationStorage) Save(o internal.Orchestration) error {
	m.Orchs[o.Name()] = o
	return nil
}

func (m OrchestrationStorage) Load(name internal.OrchestrationName) (
	internal.Orchestration,
	error,
) {
	o, exists := m.Orchs[name]
	if !exists {
		err := &internal.NotFound{Type: "orchestration", Ident: name.Unwrap()}
		return internal.Orchestration{}, err
	}
	return o, nil
}

type StageStorage struct {
	Stages map[internal.StageName]internal.Stage
}

func (m StageStorage) Save(s internal.Stage) error {
	m.Stages[s.Name()] = s
	return nil
}

func (m StageStorage) Load(name internal.StageName) (
	internal.Stage,
	error,
) {
	s, exists := m.Stages[name]
	if !exists {
		err := &internal.NotFound{Type: "stage", Ident: name.Unwrap()}
		return internal.Stage{}, err
	}
	return s, nil
}

type LinkStorage struct {
	Links map[internal.LinkName]internal.Link
}

func (m LinkStorage) Save(l internal.Link) error {
	m.Links[l.Name()] = l
	return nil
}

func (m LinkStorage) Load(name internal.LinkName) (internal.Link, error) {
	l, exists := m.Links[name]
	if !exists {
		err := &internal.NotFound{Type: "link", Ident: name.Unwrap()}
		return internal.Link{}, err
	}
	return l, nil
}