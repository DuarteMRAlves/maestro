package server

import (
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/logs"
	"github.com/dgraph-io/badger/v3"
)

func (s *Server) CreateOrchestration(req *api.CreateOrchestrationRequest) error {
	logs.LogCreateOrchestrationRequest(s.logger, req)
	return s.db.Update(
		func(txn *badger.Txn) error {
			return s.storageManager.CreateOrchestration(txn, req)
		},
	)
}

func (s *Server) GetOrchestration(
	req *api.GetOrchestrationRequest,
) ([]*api.Orchestration, error) {
	var (
		orchestrations []*api.Orchestration
		err            error
	)
	logs.LogGetOrchestrationRequest(s.logger, req)
	err = s.db.View(
		func(txn *badger.Txn) error {
			orchestrations, err = s.storageManager.GetMatchingOrchestration(
				txn,
				req,
			)
			return err
		},
	)
	if err != nil {
		return nil, err
	}
	return orchestrations, nil
}
