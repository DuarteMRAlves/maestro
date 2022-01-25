package server

import (
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/dgraph-io/badger/v3"
	"go.uber.org/zap"
)

// CreateLink creates a new link with the specified config.
// It returns an error if the asset can not be created and nil otherwise.
func (s *Server) CreateLink(cfg *apitypes.Link) error {
	s.logger.Info("Create Link.", logLink(cfg, "cfg")...)
	return s.db.Update(
		func(txn *badger.Txn) error {
			l, err := s.orchestrationManager.CreateLink(txn, cfg)
			if err != nil {
				return err
			}
			source, ok := s.orchestrationManager.GetStageByName(
				txn,
				cfg.SourceStage,
			)
			if !ok {
				return errdefs.InternalWithMsg("source not found")
			}
			target, ok := s.orchestrationManager.GetStageByName(
				txn,
				cfg.TargetStage,
			)
			if !ok {
				return errdefs.InternalWithMsg("target not found")
			}
			return s.flowManager.RegisterLink(source, target, l)
		},
	)

}

func (s *Server) GetLink(query *apitypes.Link) ([]*apitypes.Link, error) {
	var (
		links []*apitypes.Link
		err   error
	)
	s.logger.Info("Get Link.", logLink(query, "query")...)
	err = s.db.View(
		func(txn *badger.Txn) error {
			links, err = s.orchestrationManager.GetMatchingLinks(txn, query)
			return err
		},
	)
	if err != nil {
		return nil, err
	}
	return links, nil
}

func logLink(l *apitypes.Link, field string) []zap.Field {
	if l == nil {
		return []zap.Field{zap.String(field, "null")}
	}
	return []zap.Field{
		zap.String("name", string(l.Name)),
		zap.String("source-stage", string(l.SourceStage)),
		zap.String("source-field", l.SourceField),
		zap.String("target-stage", string(l.TargetStage)),
		zap.String("target-field", l.TargetField),
	}
}
