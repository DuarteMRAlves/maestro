package server

import (
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/dgraph-io/badger/v3"
	"go.uber.org/zap"
)

// CreateAsset creates a new asset with the specified config.
// It returns an error if the asset can not be created and nil otherwise.
func (s *Server) CreateAsset(req *api.CreateAssetRequest) error {
	var logFields []zap.Field
	if req == nil {
		logFields = []zap.Field{zap.String("request", "null")}
	} else {
		logFields = []zap.Field{
			zap.String("name", string(req.Name)),
			zap.String("image", req.Image),
		}
	}
	s.logger.Info("Create Asset.", logFields...)
	return s.db.Update(
		func(txn *badger.Txn) error {
			return s.storageManager.CreateAsset(txn, req)
		},
	)
}

func (s *Server) GetAsset(req *api.GetAssetRequest) (
	[]*api.Asset,
	error,
) {
	var (
		assets    []*api.Asset
		err       error
		logFields []zap.Field
	)
	if req == nil {
		logFields = []zap.Field{zap.String("request", "null")}
	} else {
		logFields = []zap.Field{
			zap.String("name", string(req.Name)),
			zap.String("image", req.Image),
		}
	}
	s.logger.Info("Get Asset.", logFields...)
	err = s.db.View(
		func(txn *badger.Txn) error {
			assets, err = s.storageManager.GetMatchingAssets(txn, req)
			return err
		},
	)
	if err != nil {
		return nil, err
	}
	return assets, nil
}
