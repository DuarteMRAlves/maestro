package server

import (
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/naming"
	"github.com/DuarteMRAlves/maestro/internal/stage"
	"github.com/DuarteMRAlves/maestro/internal/validate"
	"go.uber.org/zap"
)

// CreateStage creates a new stage with the specified config.
// It returns an error if the asset can not be created and nil otherwise.
func (s *Server) CreateStage(config *stage.Stage) error {
	s.logger.Info("Create Stage.", logStage(config, "config")...)
	if err := s.validateCreateStageConfig(config); err != nil {
		return err
	}
	return s.stageStore.Create(config)
}

func (s *Server) GetStage(query *stage.Stage) []*stage.Stage {
	s.logger.Info("Get Stage.", logStage(query, "query")...)
	return s.stageStore.Get(query)
}

func logStage(s *stage.Stage, field string) []zap.Field {
	if s == nil {
		return []zap.Field{zap.String(field, "null")}
	}
	return []zap.Field{
		zap.String("name", s.Name),
		zap.String("asset", s.Asset),
		zap.String("service", s.Service),
		zap.String("method", s.Method),
		zap.String("address", s.Address),
	}
}

// validateCreateStageConfig verifies if all conditions to create a stage are met.
// It returns an error if a condition is not met and nil otherwise.
func (s *Server) validateCreateStageConfig(config *stage.Stage) error {
	if ok, err := validate.ArgNotNil(config, "config"); !ok {
		return err
	}
	if !naming.IsValidName(config.Name) {
		return errdefs.InvalidArgumentWithMsg(
			"invalid name '%v'",
			config.Name)
	}
	// Asset is not required but if specified should exist.
	if config.Asset != "" && !s.assetStore.Contains(config.Asset) {
		return errdefs.NotFoundWithMsg(
			"asset '%v' not found",
			config.Asset)
	}
	return nil
}
