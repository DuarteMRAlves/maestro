package server

import (
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/link"
	"github.com/DuarteMRAlves/maestro/internal/naming"
	"github.com/DuarteMRAlves/maestro/internal/validate"
	"go.uber.org/zap"
)

// CreateLink creates a new link with the specified config.
// It returns an error if the asset can not be created and nil otherwise.
func (s *Server) CreateLink(config *link.Link) error {
	s.logger.Info("Create Link.", logLink(config, "config")...)
	if err := s.validateCreateLinkConfig(config); err != nil {
		return err
	}
	return s.linkStore.Create(config)
}

func (s *Server) GetLink(query *link.Link) []*link.Link {
	s.logger.Info("Get Link.", logLink(query, "query")...)
	return s.linkStore.Get(query)
}

func logLink(l *link.Link, field string) []zap.Field {
	if l == nil {
		return []zap.Field{zap.String(field, "null")}
	}
	return []zap.Field{
		zap.String("name", l.Name),
		zap.String("source-stage", l.SourceStage),
		zap.String("source-field", l.SourceField),
		zap.String("target-stage", l.TargetStage),
		zap.String("target-field", l.TargetField),
	}
}

// validateCreateLinkConfig verifies if all conditions to create a link are met.
// It returns an error if a condition is not met and nil otherwise.
func (s *Server) validateCreateLinkConfig(config *link.Link) error {
	if ok, err := validate.ArgNotNil(config, "config"); !ok {
		return err
	}
	if !naming.IsValidName(config.Name) {
		return errdefs.InvalidArgumentWithMsg("invalid name '%v'", config.Name)
	}
	if config.SourceStage == "" {
		return errdefs.InvalidArgumentWithMsg("empty source stage name")
	}
	if config.TargetStage == "" {
		return errdefs.InvalidArgumentWithMsg("empty target stage name")
	}
	if config.SourceStage == config.TargetStage {
		return errdefs.InvalidArgumentWithMsg(
			"source and target stages are equal")
	}
	source, ok := s.stageStore.GetByName(config.SourceStage)
	if !ok {
		return errdefs.NotFoundWithMsg(
			"source stage '%v' not found",
			config.SourceStage)
	}
	target, ok := s.stageStore.GetByName(config.TargetStage)
	if !ok {
		return errdefs.NotFoundWithMsg(
			"target stage '%v' not found",
			config.TargetStage)
	}

	sourceOutput := source.Rpc.Output()
	targetInput := target.Rpc.Input()
	if config.SourceField != "" {
		sourceOutput, ok = sourceOutput.GetMessageField(config.SourceField)
		if !ok {
			return errdefs.NotFoundWithMsg(
				"field with name %s not found for message %s for source stage "+
					"in link %s",
				config.SourceField,
				source.Rpc.Output().FullyQualifiedName(),
				config.Name)
		}
	}
	if config.TargetField != "" {
		targetInput, ok = targetInput.GetMessageField(config.TargetField)
		if !ok {
			return errdefs.NotFoundWithMsg(
				"field with name %s not found for message %s for target stage "+
					"in link %v",
				config.TargetField,
				target.Rpc.Input().FullyQualifiedName(),
				config.Name)
		}
	}
	if !sourceOutput.Compatible(targetInput) {
		return errdefs.InvalidArgumentWithMsg(
			"incompatible message types between source output %s and target"+
				" input %s in link %s",
			sourceOutput.FullyQualifiedName(),
			targetInput.FullyQualifiedName(),
			config.Name)
	}
	return nil
}
