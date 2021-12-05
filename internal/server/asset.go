package server

import (
	"github.com/DuarteMRAlves/maestro/internal/assert"
	"github.com/DuarteMRAlves/maestro/internal/asset"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/naming"
	"log"
)

// CreateAsset creates a new asset with the specified config.
// It returns an error if the asset can not be created and nil otherwise.
func (s *Server) CreateAsset(config *asset.Asset) error {
	log.Printf("Create Asset with config='%v'\n", config)
	if err := s.validateCreateAssetConfig(config); err != nil {
		return err
	}
	return s.assetStore.Create(config)
}

func (s *Server) GetAsset(query *asset.Asset) []*asset.Asset {
	log.Printf("Get Asset with query='%v'\n", query)
	return s.assetStore.Get(query)
}

// validateCreateAssetConfig verifies if all conditions to create an asset are met.
// It returns an error if a condition is not met and nil otherwise.
func (s *Server) validateCreateAssetConfig(config *asset.Asset) error {
	if ok, err := assert.ArgNotNil(config, "config"); !ok {
		return errdefs.InvalidArgumentWithError(err)
	}
	if !naming.IsValidName(config.Name) {
		return errdefs.InvalidArgumentWithMsg(
			"invalid name '%v'",
			config.Name)
	}
	return nil
}
