package create

import "github.com/DuarteMRAlves/maestro/internal/domain"

type SaveAsset func(domain.Asset) domain.AssetResult
type LoadAsset func(domain.AssetName) domain.AssetResult
type ExistsAsset func(domain.AssetName) bool

type SaveStage func(domain.Stage) domain.StageResult
type LoadStage func(domain.StageName) domain.StageResult

type SaveLink func(domain.Link) domain.LinkResult
type LoadLink func(domain.LinkName) domain.LinkResult

type SaveOrchestration func(domain.Orchestration) domain.OrchestrationResult
type LoadOrchestration func(domain.OrchestrationName) domain.OrchestrationResult

type AssetRequest struct {
	Name  string
	Image domain.OptionalString
}

type AssetResponse struct {
	Err domain.OptionalError
}
