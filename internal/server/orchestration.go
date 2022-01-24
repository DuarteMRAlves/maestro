package server

import (
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"github.com/dgraph-io/badger/v3"
	"go.uber.org/zap"
)

func (s *Server) CreateOrchestration(config *apitypes.Orchestration) error {
	s.logger.Info(
		"Create Orchestration.",
		logOrchestration(config, "config")...)
	return s.db.Update(func(txn *badger.Txn) error {
		return s.orchestrationManager.CreateOrchestration(txn, config)
	})
}

func (s *Server) GetOrchestration(
	query *apitypes.Orchestration,
) ([]*apitypes.Orchestration, error) {
	var (
		orchestrations []*apitypes.Orchestration
		err            error
	)
	s.logger.Info("Get Orchestration.", logOrchestration(query, "query")...)
	err = s.db.View(func(txn *badger.Txn) error {
		orchestrations, err = s.orchestrationManager.GetMatchingOrchestration(
			txn,
			query)
		return err
	})
	if err != nil {
		return nil, err
	}
	return orchestrations, nil
}

func logOrchestration(
	o *apitypes.Orchestration,
	field string,
) []zap.Field {
	if o == nil {
		return []zap.Field{zap.String(field, "null")}
	}
	links := make([]string, 0, len(o.Links))
	for _, l := range o.Links {
		links = append(links, string(l))
	}
	return []zap.Field{
		zap.String("name", string(o.Name)),
		zap.Strings("links", links),
	}
}
