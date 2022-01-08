package server

import (
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/link"
	"github.com/DuarteMRAlves/maestro/internal/naming"
	"github.com/DuarteMRAlves/maestro/internal/validate"
	"go.uber.org/zap"
)

// CreateLink creates a new link with the specified config.
// It returns an error if the asset can not be created and nil otherwise.
func (s *Server) CreateLink(config *apitypes.Link) error {
	s.logger.Info("Create Link.", logLink(config, "config")...)
	if err := s.validateCreateLinkConfig(config); err != nil {
		return err
	}
	l := link.New(
		config.Name,
		config.SourceStage,
		config.SourceField,
		config.TargetStage,
		config.TargetField,
	)
	source, ok := s.stageStore.GetByName(config.SourceStage)
	if !ok {
		return errdefs.InternalWithMsg("source not found")
	}
	target, ok := s.stageStore.GetByName(config.TargetStage)
	if !ok {
		return errdefs.InternalWithMsg("target not found")
	}
	if err := s.flowManager.Register(source, target, l); err != nil {
		return err
	}
	return s.linkStore.Create(l)
}

func (s *Server) GetLink(query *apitypes.Link) []*apitypes.Link {
	s.logger.Info("Get Link.", logLink(query, "query")...)
	links := s.linkStore.Get(query)
	apiLinks := make([]*apitypes.Link, 0, len(links))
	for _, l := range links {
		apiLinks = append(apiLinks, l.ToApi())
	}
	return apiLinks
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

// validateCreateLinkConfig verifies if all conditions to create a link are met.
// It returns an error if a condition is not met and nil otherwise.
func (s *Server) validateCreateLinkConfig(config *apitypes.Link) error {
	if ok, err := validate.ArgNotNil(config, "config"); !ok {
		return err
	}
	if !naming.IsValidLinkName(config.Name) {
		return errdefs.InvalidArgumentWithMsg("invalid name '%v'", config.Name)
	}
	if s.linkStore.Contains(config.Name) {
		return errdefs.AlreadyExistsWithMsg("link '%v' already exists", config.Name)
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

	if !source.IsPending() {
		return errdefs.FailedPreconditionWithMsg(
			"source stage is not in Pending phase for link %s",
			config.Name)
	}
	if !target.IsPending() {
		return errdefs.FailedPreconditionWithMsg(
			"target stage is not in Pending phase for link %s",
			config.Name)
	}
	return nil
}
