package storage

import apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"

func copyOrchestration(
	dst *apitypes.Orchestration,
	src *apitypes.Orchestration,
) {
	dst.Name = src.Name
	dst.Phase = src.Phase
	dst.Stages = make([]apitypes.StageName, 0, len(src.Stages))
	for _, s := range src.Stages {
		dst.Stages = append(dst.Stages, s)
	}
	dst.Links = make([]apitypes.LinkName, 0, len(src.Links))
	for _, l := range src.Links {
		dst.Links = append(dst.Links, l)
	}
}
