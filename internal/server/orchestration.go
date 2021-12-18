package server

import (
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/naming"
	"github.com/DuarteMRAlves/maestro/internal/orchestration"
	"github.com/DuarteMRAlves/maestro/internal/validate"
	"go.uber.org/zap"
)

// CreateOrchestration creates a orchestration from the given config.
// The function returns an error if the orchestration name is not valid or if
// one of the links does not exist.
func (s *Server) CreateOrchestration(config *orchestration.Orchestration) error {
	s.logger.Info(
		"Create Orchestration.",
		logOrchestration(config, "config")...)
	if err := s.validateCreateOrchestrationConfig(config); err != nil {
		return err
	}
	return s.orchestrationStore.Create(config)
}

// GetOrchestration returns a list of orchestrations that match the received query.
func (s *Server) GetOrchestration(
	query *orchestration.Orchestration,
) []*orchestration.Orchestration {
	s.logger.Info("Get Orchestration.", logOrchestration(query, "query")...)
	return s.orchestrationStore.Get(query)
}

func logOrchestration(
	o *orchestration.Orchestration,
	field string,
) []zap.Field {
	if o == nil {
		return []zap.Field{zap.String(field, "null")}
	}
	return []zap.Field{
		zap.String("name", o.Name),
		zap.Strings("links", o.Links),
	}
}

// validateCreateOrchestrationConfig verifies if all the conditions to create a
// orchestration are met. It returns an error if one condition is not met and nil
// otherwise.
func (s *Server) validateCreateOrchestrationConfig(
	config *orchestration.Orchestration,
) error {
	if ok, err := validate.ArgNotNil(config, "config"); !ok {
		return err
	}
	if !naming.IsValidName(config.Name) {
		return errdefs.InvalidArgumentWithMsg("invalid name '%v'", config.Name)
	}
	for _, l := range config.Links {
		if !s.linkStore.Contains(l) {
			return errdefs.NotFoundWithMsg("link '%v' not found", l)
		}
	}
	return nil
}
