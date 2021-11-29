package server

import (
	"github.com/DuarteMRAlves/maestro/internal/blueprint"
	"log"
)

func (s *Server) CreateBlueprint(config *blueprint.Blueprint) error {
	log.Printf("Create Blueprint with config=%v", config)
	return s.blueprintStore.Create(config)
}
