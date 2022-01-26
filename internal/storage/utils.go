package storage

import (
	"github.com/DuarteMRAlves/maestro/internal/api"
)

func copyOrchestration(
	dst *api.Orchestration,
	src *api.Orchestration,
) {
	dst.Name = src.Name
	dst.Phase = src.Phase
	dst.Stages = make([]api.StageName, 0, len(src.Stages))
	for _, s := range src.Stages {
		dst.Stages = append(dst.Stages, s)
	}
	dst.Links = make([]api.LinkName, 0, len(src.Links))
	for _, l := range src.Links {
		dst.Links = append(dst.Links, l)
	}
}
