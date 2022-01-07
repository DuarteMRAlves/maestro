package server

import (
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/naming"
	"github.com/DuarteMRAlves/maestro/internal/orchestration"
	"github.com/DuarteMRAlves/maestro/internal/validate"
	"go.uber.org/zap"
)

// CreateOrchestration creates a orchestration from the given config.
// The function returns an error if the orchestration name is not valid or if
// one of the links does not exist.
func (s *Server) CreateOrchestration(config *apitypes.Orchestration) error {
	s.logger.Info(
		"Create Orchestration.",
		logOrchestration(config, "config")...)
	if err := s.validateCreateOrchestrationConfig(config); err != nil {
		return err
	}
	o := orchestration.New(config.Name, config.Links)
	return s.orchestrationStore.Create(o)
}

// GetOrchestration returns a list of orchestrations that match the received query.
func (s *Server) GetOrchestration(
	query *apitypes.Orchestration,
) []*apitypes.Orchestration {
	s.logger.Info("Get Orchestration.", logOrchestration(query, "query")...)
	orchestrations := s.orchestrationStore.Get(query)
	apiOrchestrations := make([]*apitypes.Orchestration, 0, len(orchestrations))
	for _, o := range orchestrations {
		apiOrchestrations = append(apiOrchestrations, o.ToApi())
	}
	return apiOrchestrations
}

func logOrchestration(
	o *apitypes.Orchestration,
	field string,
) []zap.Field {
	if o == nil {
		return []zap.Field{zap.String(field, "null")}
	}
	return []zap.Field{
		zap.String("name", string(o.Name)),
		zap.Strings("links", o.Links),
	}
}

// validateCreateOrchestrationConfig verifies if all the conditions to create a
// orchestration are met. It returns an error if one condition is not met and nil
// otherwise.
func (s *Server) validateCreateOrchestrationConfig(
	config *apitypes.Orchestration,
) error {
	if ok, err := validate.ArgNotNil(config, "config"); !ok {
		return err
	}
	if !naming.IsValidName(string(config.Name)) {
		return errdefs.InvalidArgumentWithMsg("invalid name '%v'", config.Name)
	}
	for _, l := range config.Links {
		if !s.linkStore.Contains(l) {
			return errdefs.NotFoundWithMsg("link '%v' not found", l)
		}
	}
	return nil
}
