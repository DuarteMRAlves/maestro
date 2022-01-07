package server

import (
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"github.com/DuarteMRAlves/maestro/internal/asset"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/naming"
	"github.com/DuarteMRAlves/maestro/internal/validate"
	"go.uber.org/zap"
)

// CreateAsset creates a new asset with the specified config.
// It returns an error if the asset can not be created and nil otherwise.
func (s *Server) CreateAsset(config *apitypes.Asset) error {
	s.logger.Info("Create Asset.", logAsset(config, "config")...)
	if err := s.validateCreateAssetConfig(config); err != nil {
		return err
	}
	a := asset.New(config.Name, config.Image)
	return s.assetStore.Create(a)
}

func (s *Server) GetAsset(query *apitypes.Asset) []*apitypes.Asset {
	s.logger.Info("Get Asset.", logAsset(query, "query")...)
	assets := s.assetStore.Get(query)
	apiAssets := make([]*apitypes.Asset, 0, len(assets))
	for _, a := range assets {
		apiAssets = append(apiAssets, a.ToApi())
	}
	return apiAssets
}

func logAsset(a *apitypes.Asset, field string) []zap.Field {
	if a == nil {
		return []zap.Field{zap.String(field, "null")}
	}
	return []zap.Field{
		zap.String("name", string(a.Name)),
		zap.String("image", a.Image),
	}
}

// validateCreateAssetConfig verifies if all conditions to create an asset are met.
// It returns an error if a condition is not met and nil otherwise.
func (s *Server) validateCreateAssetConfig(config *apitypes.Asset) error {
	if ok, err := validate.ArgNotNil(config, "config"); !ok {
		return errdefs.InvalidArgumentWithError(err)
	}
	if !naming.IsValidAssetName(config.Name) {
		return errdefs.InvalidArgumentWithMsg(
			"invalid name '%v'",
			config.Name)
	}
	return nil
}
