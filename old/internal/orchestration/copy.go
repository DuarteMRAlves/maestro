package orchestration

import (
	"github.com/DuarteMRAlves/maestro/internal/api"
)

func copyOrchestration(dst *api.Orchestration, src *api.Orchestration) {
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

func copyStage(dst *api.Stage, src *api.Stage) {
	dst.Name = src.Name
	dst.Phase = src.Phase
	dst.Service = src.Service
	dst.Rpc = src.Rpc
	dst.Address = src.Address
	dst.Orchestration = src.Orchestration
	dst.Asset = src.Asset
}

func copyLink(dst *api.Link, src *api.Link) {
	dst.Name = src.Name
	dst.SourceStage = src.SourceStage
	dst.SourceField = src.SourceField
	dst.TargetStage = src.TargetStage
	dst.TargetField = src.TargetField
	dst.Orchestration = src.Orchestration
}

func copyAsset(dst *api.Asset, src *api.Asset) {
	dst.Name = src.Name
	dst.Image = src.Image
}
