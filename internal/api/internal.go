package api

import (
	"github.com/DuarteMRAlves/maestro/internal/asset"
	"github.com/DuarteMRAlves/maestro/internal/blueprint"
	"log"
)

// InternalAPI is an interface that collects all the available commands
// for the maestro server. All calls on external APIs should be redirected
// through this API that collects all functionality.
type InternalAPI interface {
	CreateAsset(config *asset.Asset) error
	GetAsset(query *asset.Asset) []*asset.Asset

	CreateBlueprint(config *blueprint.Blueprint) error
}

type internalAPI struct {
	assetStore     asset.Store
	blueprintStore blueprint.Store
}

func NewInternalAPI() InternalAPI {
	return &internalAPI{
		assetStore:     asset.NewStore(),
		blueprintStore: blueprint.NewStore(),
	}
}

func (m *internalAPI) CreateAsset(config *asset.Asset) error {
	log.Printf("Create Asset with config='%v'\n", config)
	return m.assetStore.Create(config)
}

func (m *internalAPI) GetAsset(query *asset.Asset) []*asset.Asset {
	log.Printf("Get Asset with query='%v'\n", query)
	return m.assetStore.Get(query)
}

func (m *internalAPI) CreateBlueprint(config *blueprint.Blueprint) error {
	log.Printf("Create Blueprint with config=%v", config)
	return m.blueprintStore.Create(config)
}
