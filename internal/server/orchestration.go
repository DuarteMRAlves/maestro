package server

import (
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/dgraph-io/badger/v3"
	"go.uber.org/zap"
)

func (s *Server) CreateOrchestration(req *api.CreateOrchestrationRequest) error {
	var logFields []zap.Field
	if req == nil {
		logFields = []zap.Field{zap.String("req", "nil")}
	} else {
		links := make([]string, 0, len(req.Links))
		for _, l := range req.Links {
			links = append(links, string(l))
		}
		logFields = []zap.Field{
			zap.String("name", string(req.Name)),
			zap.Strings("links", links),
		}
	}
	s.logger.Info("Create Orchestration.", logFields...)
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
		logFields      []zap.Field
	)
	if req == nil {
		logFields = []zap.Field{zap.String("req", "nil")}
	} else {
		logFields = []zap.Field{zap.String("name", string(req.Name))}
	}
	s.logger.Info("Get Orchestration.", logFields...)
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
