package server

import (
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"github.com/DuarteMRAlves/maestro/internal/asset"
	"github.com/dgraph-io/badger/v3"
	"go.uber.org/zap"
)

// CreateAsset creates a new asset with the specified config.
// It returns an error if the asset can not be created and nil otherwise.
func (s *Server) CreateAsset(cfg *apitypes.Asset) error {
	s.logger.Info("Create Asset.", logAsset(cfg, "cfg")...)
	return s.db.Update(func(txn *badger.Txn) error {
		return asset.Create(txn, cfg)
	})
}

func (s *Server) GetAsset(query *apitypes.Asset) ([]*apitypes.Asset, error) {
	var (
		assets []*apitypes.Asset
		err    error
	)
	s.logger.Info("Get Asset.", logAsset(query, "query")...)
	err = s.db.View(func(txn *badger.Txn) error {
		assets, err = asset.Get(txn, query)
		return err
	})
	if err != nil {
		return nil, err
	}
	return assets, nil
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
