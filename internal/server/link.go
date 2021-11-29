package server

import (
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/link"
	"log"
)

func (s *Server) CreateLink(config *link.Link) error {
	log.Printf("Create Stage with config='%v'\n", config)
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
	return s.linkStore.Create(config)
}

func (s *Server) GetLink(query *link.Link) []*link.Link {
	log.Printf("Get Link with query=%v", query)
	return s.linkStore.Get(query)
}
