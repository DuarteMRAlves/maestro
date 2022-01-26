package server

import (
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/dgraph-io/badger/v3"
	"go.uber.org/zap"
)

// CreateLink creates a new link with the specified config.
// It returns an error if the asset can not be created and nil otherwise.
func (s *Server) CreateLink(req *api.CreateLinkRequest) error {
	var logFields []zap.Field
	if req == nil {
		logFields = []zap.Field{zap.String("req", "null")}
	} else {
		logFields = []zap.Field{
			zap.String("name", string(req.Name)),
			zap.String("source-stage", string(req.SourceStage)),
			zap.String("source-field", req.SourceField),
			zap.String("target-stage", string(req.TargetStage)),
			zap.String("target-field", req.TargetField),
		}
	}
	s.logger.Info("Create Link.", logFields...)
	return s.db.Update(
		func(txn *badger.Txn) error {
			l, err := s.storageManager.CreateLink(txn, req)
			if err != nil {
				return err
			}
			source, ok := s.storageManager.GetStageByName(
				txn,
				req.SourceStage,
			)
			if !ok {
				return errdefs.InternalWithMsg("source not found")
			}
			target, ok := s.storageManager.GetStageByName(
				txn,
				req.TargetStage,
			)
			if !ok {
				return errdefs.InternalWithMsg("target not found")
			}
			return s.flowManager.RegisterLink(source, target, l)
		},
	)

}

func (s *Server) GetLink(req *api.GetLinkRequest) ([]*api.Link, error) {
	var (
		links     []*api.Link
		err       error
		logFields []zap.Field
	)
	if req == nil {
		logFields = []zap.Field{zap.String("req", "null")}
	} else {
		logFields = []zap.Field{
			zap.String("name", string(req.Name)),
			zap.String("source-stage", string(req.SourceStage)),
			zap.String("source-field", req.SourceField),
			zap.String("target-stage", string(req.TargetStage)),
			zap.String("target-field", req.TargetField),
		}
	}
	s.logger.Info("Get Link.", logFields...)
	err = s.db.View(
		func(txn *badger.Txn) error {
			links, err = s.storageManager.GetMatchingLinks(txn, req)
			return err
		},
	)
	if err != nil {
		return nil, err
	}
	return links, nil
}
