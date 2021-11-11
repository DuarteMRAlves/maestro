package api

import (
	"github.com/DuarteMRAlves/maestro/internal/asset"
	"github.com/DuarteMRAlves/maestro/internal/identifier"
	"log"
)

// InternalAPI is an interface that collects all the available commands
// for the maestro server. All calls on external APIs should be redirected
// through this API that collects all functionality.
type InternalAPI interface {
	asset.Store
}

type internalAPI struct {
	assetStore asset.Store
}

func NewInternalAPI() InternalAPI {
	return &internalAPI{
		assetStore: asset.NewStore(),
	}
}

func (m *internalAPI) Create(description *asset.Asset) (identifier.Id, error) {
	log.Printf("Create request with description='%v'\n", description)
	return m.assetStore.Create(description)
}

func (m *internalAPI) Get(id identifier.Id) (*asset.Asset, error) {
	log.Printf("Get request with identifier='%v'\n", id)
	return m.assetStore.Get(id)
}

func (m *internalAPI) List() ([]*asset.Asset, error) {
	log.Printf("List request")
	return m.assetStore.List()
}
