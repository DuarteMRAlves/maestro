package server

import (
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"github.com/dgraph-io/badger/v3"
	"go.uber.org/zap"
)

// CreateStage creates a new stage with the specified config.
// It returns an error if the asset can not be created and nil otherwise.
func (s *Server) CreateStage(cfg *apitypes.Stage) error {
	s.logger.Info("Create Stage.", logStage(cfg, "cfg")...)
	return s.db.Update(
		func(txn *badger.Txn) error {
			st, err := s.storageManager.CreateStage(txn, cfg)
			if err != nil {
				return err
			}
			return s.flowManager.RegisterStage(st)
		},
	)
}

func (s *Server) GetStage(query *apitypes.Stage) ([]*apitypes.Stage, error) {
	var (
		stages []*apitypes.Stage
		err    error
	)
	s.logger.Info("Get Stage.", logStage(query, "query")...)
	err = s.db.View(
		func(txn *badger.Txn) error {
			stages, err = s.storageManager.GetMatchingStage(txn, query)
			return err
		},
	)
	if err != nil {
		return nil, err
	}
	return stages, nil
}

func logStage(s *apitypes.Stage, field string) []zap.Field {
	if s == nil {
		return []zap.Field{zap.String(field, "null")}
	}
	return []zap.Field{
		zap.String("name", string(s.Name)),
		zap.String("asset", string(s.Asset)),
		zap.String("service", s.Service),
		zap.String("rpc", s.Rpc),
		zap.String("address", s.Address),
		zap.String("host", s.Host),
		zap.Int32("port", s.Port),
	}
}
