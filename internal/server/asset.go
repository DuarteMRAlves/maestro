package server

import (
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/logs"
	"github.com/dgraph-io/badger/v3"
)

// CreateAsset creates a new asset with the specified config.
// It returns an error if the asset can not be created and nil otherwise.
func (s *Server) CreateAsset(req *api.CreateAssetRequest) error {
	logs.LogCreateAssetRequest(s.logger, req)
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
		assets []*api.Asset
		err    error
	)
	logs.LogGetAssetRequest(s.logger, req)
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
