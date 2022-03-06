package create

import (
	"github.com/DuarteMRAlves/maestro/internal/domain"
	"github.com/DuarteMRAlves/maestro/internal/execute"
)

type SaveAsset func(domain.Asset) domain.AssetResult
type LoadAsset func(domain.AssetName) domain.AssetResult
type ExistsAsset func(domain.AssetName) bool

type AssetRequest struct {
	Name  string
	Image domain.OptionalString
}

type AssetResponse struct {
	Err domain.OptionalError
}

type SaveStage func(domain.Stage) domain.StageResult
type LoadStage func(domain.StageName) domain.StageResult

type SaveLink func(execute.Link) execute.LinkResult
type LoadLink func(domain.LinkName) execute.LinkResult

type SaveOrchestration func(Orchestration) OrchestrationResult
type ExistsOrchestration func(domain.OrchestrationName) bool

type Orchestration interface {
	Name() domain.OrchestrationName
	Stages() []domain.StageName
	Links() []domain.LinkName
}

type OrchestrationRequest struct {
	Name string
}

type OrchestrationResponse struct {
	Err domain.OptionalError
}
