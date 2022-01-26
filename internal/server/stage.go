package server

import (
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/logs"
	"github.com/dgraph-io/badger/v3"
)

// CreateStage creates a new stage with the specified config.
// It returns an error if the asset can not be created and nil otherwise.
func (s *Server) CreateStage(req *api.CreateStageRequest) error {
	logs.LogCreateStageRequest(s.logger, req)
	return s.db.Update(
		func(txn *badger.Txn) error {
			return s.storageManager.CreateStage(txn, req)
		},
	)
}

func (s *Server) GetStage(req *api.GetStageRequest) ([]*api.Stage, error) {
	var (
		stages []*api.Stage
		err    error
	)
	logs.LogGetStageRequest(s.logger, req)
	err = s.db.View(
		func(txn *badger.Txn) error {
			stages, err = s.storageManager.GetMatchingStage(txn, req)
			return err
		},
	)
	if err != nil {
		return nil, err
	}
	return stages, nil
}
