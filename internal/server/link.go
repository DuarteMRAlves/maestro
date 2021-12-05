package server

import (
	"github.com/DuarteMRAlves/maestro/internal/assert"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/link"
	"github.com/DuarteMRAlves/maestro/internal/naming"
	"log"
)

// CreateLink creates a new link with the specified config.
// It returns an error if the asset can not be created and nil otherwise.
func (s *Server) CreateLink(config *link.Link) error {
	log.Printf("Create Stage with config='%v'\n", config)
	if err := s.validateCreateLinkConfig(config); err != nil {
		return err
	}
	return s.linkStore.Create(config)
}

func (s *Server) GetLink(query *link.Link) []*link.Link {
	log.Printf("Get Link with query=%v", query)
	return s.linkStore.Get(query)
}

// validateCreateLinkConfig verifies if all conditions to create a link are met.
// It returns an error if a condition is not met and nil otherwise.
func (s *Server) validateCreateLinkConfig(config *link.Link) error {
	if ok, err := assert.ArgNotNil(config, "config"); !ok {
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
	if !s.stageStore.Contains(config.SourceStage) {
		return errdefs.NotFoundWithMsg(
			"source stage '%v' not found",
			config.SourceStage)
	}
	if !s.stageStore.Contains(config.TargetStage) {
		return errdefs.NotFoundWithMsg(
			"target stage '%v' not found",
			config.TargetStage)
	}
	return nil
}
