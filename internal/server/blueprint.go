package server

import (
	"github.com/DuarteMRAlves/maestro/internal/blueprint"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/naming"
	"github.com/DuarteMRAlves/maestro/internal/validate"
	"log"
)

// CreateBlueprint creates a blueprint from the given config.
// The function returns an error if the blueprint name is not valid or if
// one of the links does not exist.
func (s *Server) CreateBlueprint(config *blueprint.Blueprint) error {
	log.Printf("Create Blueprint with config=%v", config)
	if err := s.validateCreateBlueprintConfig(config); err != nil {
		return err
	}
	return s.blueprintStore.Create(config)
}

// GetBlueprint returns a list of blueprints that match the received query.
func (s *Server) GetBlueprint(
	query *blueprint.Blueprint,
) []*blueprint.Blueprint {
	log.Printf("Get Blueprint with query='%v'\n", query)
	return s.blueprintStore.Get(query)
}

// validateCreateBlueprintConfig verifies if all the conditions to create a
// blueprint are met. It returns an error if one condition is not met and nil
// otherwise.
func (s *Server) validateCreateBlueprintConfig(
	config *blueprint.Blueprint,
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
