package server

import (
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/dgraph-io/badger/v3"
	"go.uber.org/zap"
)

// CreateStage creates a new stage with the specified config.
// It returns an error if the asset can not be created and nil otherwise.
func (s *Server) CreateStage(req *api.CreateStageRequest) error {
	var logFields []zap.Field
	if req == nil {
		logFields = []zap.Field{zap.String("req", "nil")}
	} else {
		logFields = []zap.Field{
			zap.String("name", string(req.Name)),
			zap.String("asset", string(req.Asset)),
			zap.String("service", req.Service),
			zap.String("rpc", req.Rpc),
			zap.String("address", req.Address),
			zap.String("host", req.Host),
			zap.Int32("port", req.Port),
		}
	}
	s.logger.Info("Create Stage.", logFields...)
	return s.db.Update(
		func(txn *badger.Txn) error {
			st, err := s.storageManager.CreateStage(txn, req)
			if err != nil {
				return err
			}
			return s.flowManager.RegisterStage(st)
		},
	)
}

func (s *Server) GetStage(req *api.GetStageRequest) ([]*api.Stage, error) {
	var (
		stages    []*api.Stage
		err       error
		logFields []zap.Field
	)
	if req == nil {
		logFields = []zap.Field{zap.String("req", "nil")}
	} else {
		logFields = []zap.Field{
			zap.String("name", string(req.Name)),
			zap.String("phase", string(req.Phase)),
			zap.String("asset", string(req.Asset)),
			zap.String("service", req.Service),
			zap.String("rpc", req.Rpc),
			zap.String("address", req.Address),
		}
	}
	s.logger.Info("Get Stage.", logFields...)
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
