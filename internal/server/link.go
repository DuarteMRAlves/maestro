package server

import (
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/dgraph-io/badger/v3"
	"go.uber.org/zap"
)

// CreateLink creates a new link with the specified config.
// It returns an error if the asset can not be created and nil otherwise.
func (s *Server) CreateLink(config *apitypes.Link) error {
	s.logger.Info("Create Link.", logLink(config, "config")...)
	return s.db.Update(
		func(txn *badger.Txn) error {
			l, err := s.orchestrationManager.CreateLink(txn, config)
			if err != nil {
				return err
			}
			source, ok := s.orchestrationManager.GetStageByName(
				txn,
				config.SourceStage,
			)
			if !ok {
				return errdefs.InternalWithMsg("source not found")
			}
			target, ok := s.orchestrationManager.GetStageByName(
				txn,
				config.TargetStage,
			)
			if !ok {
				return errdefs.InternalWithMsg("target not found")
			}
			return s.flowManager.RegisterLink(source, target, l)
		},
	)

}

func (s *Server) GetLink(query *apitypes.Link) []*apitypes.Link {
	s.logger.Info("Get Link.", logLink(query, "query")...)
	return s.orchestrationManager.GetMatchingLinks(query)
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
