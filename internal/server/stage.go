package server

import (
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/naming"
	"github.com/DuarteMRAlves/maestro/internal/stage"
	"github.com/DuarteMRAlves/maestro/internal/validate"
	"log"
)

// CreateStage creates a new stage with the specified config.
// It returns an error if the asset can not be created and nil otherwise.
func (s *Server) CreateStage(config *stage.Stage) error {
	log.Printf("Create Stage with config='%v'\n", config)
	if err := s.validateCreateStageConfig(config); err != nil {
		return err
	}
	return s.stageStore.Create(config)
}

func (s *Server) GetStage(query *stage.Stage) []*stage.Stage {
	log.Printf("Get Stage with query=%v", query)
	return s.stageStore.Get(query)
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
	if config.Asset == "" {
		return errdefs.InvalidArgumentWithMsg("empty asset name")
	}
	if !s.assetStore.Contains(config.Asset) {
		return errdefs.NotFoundWithMsg(
			"asset '%v' not found",
			config.Asset)
	}
	return nil
}
