package server

import (
	"github.com/DuarteMRAlves/maestro/internal/asset"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/naming"
	"github.com/DuarteMRAlves/maestro/internal/validate"
	"go.uber.org/zap"
)

// CreateAsset creates a new asset with the specified config.
// It returns an error if the asset can not be created and nil otherwise.
func (s *Server) CreateAsset(config *asset.Asset) error {
	s.logger.Info("Create Asset.", logAsset(config, "config")...)
	if err := s.validateCreateAssetConfig(config); err != nil {
		return err
	}
	return s.assetStore.Create(config)
}

func (s *Server) GetAsset(query *asset.Asset) []*asset.Asset {
	s.logger.Info("Get Asset.", logAsset(query, "query")...)
	return s.assetStore.Get(query)
}

func logAsset(a *asset.Asset, field string) []zap.Field {
	if a == nil {
		return []zap.Field{zap.String(field, "null")}
	}
	return []zap.Field{zap.String("name", a.Name), zap.String("image", a.Image)}
}

// validateCreateAssetConfig verifies if all conditions to create an asset are met.
// It returns an error if a condition is not met and nil otherwise.
func (s *Server) validateCreateAssetConfig(config *asset.Asset) error {
	if ok, err := validate.ArgNotNil(config, "config"); !ok {
		return errdefs.InvalidArgumentWithError(err)
	}
	if !naming.IsValidName(config.Name) {
		return errdefs.InvalidArgumentWithMsg(
			"invalid name '%v'",
			config.Name)
	}
	return nil
}
