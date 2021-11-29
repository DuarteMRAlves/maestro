package server

import (
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/stage"
	"log"
)

func (s *Server) CreateStage(config *stage.Stage) error {
	log.Printf("Create Stage with config='%v'\n", config)
	if config.Asset == "" {
		return errdefs.InvalidArgumentWithMsg("empty asset name")
	}
	if !s.assetStore.Contains(config.Asset) {
		return errdefs.NotFoundWithMsg("asset '%v' not found", config.Asset)
	}
	return s.stageStore.Create(config)
}

func (s *Server) GetStage(query *stage.Stage) []*stage.Stage {
	log.Printf("Get Stage with query=%v", query)
	return s.stageStore.Get(query)
}
