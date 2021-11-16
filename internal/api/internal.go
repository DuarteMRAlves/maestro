package api

import (
	"github.com/DuarteMRAlves/maestro/internal/asset"
	"github.com/DuarteMRAlves/maestro/internal/blueprint"
	"github.com/DuarteMRAlves/maestro/internal/identifier"
	"log"
)

// InternalAPI is an interface that collects all the available commands
// for the maestro server. All calls on external APIs should be redirected
// through this API that collects all functionality.
type InternalAPI interface {
	CreateAsset(config *asset.Asset) (identifier.Id, error)
	GetAsset(id identifier.Id) (*asset.Asset, error)
	ListAssets() ([]*asset.Asset, error)

	CreateBlueprint(config *blueprint.Blueprint) (identifier.Id, error)
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

func (m *internalAPI) CreateAsset(config *asset.Asset) (identifier.Id, error) {
	log.Printf("Create Asset with config='%v'\n", config)
	return m.assetStore.Create(config)
}

func (m *internalAPI) GetAsset(id identifier.Id) (*asset.Asset, error) {
	log.Printf("Get Asset with identifier='%v'\n", id)
	return m.assetStore.Get(id)
}

func (m *internalAPI) ListAssets() ([]*asset.Asset, error) {
	log.Printf("List Assets")
	return m.assetStore.List()
}

func (m *internalAPI) CreateBlueprint(
	config *blueprint.Blueprint,
) (identifier.Id, error) {
	log.Printf("Create Blueprint with config=%v", config)
	return m.blueprintStore.Create(config)
}
