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

type SaveStage func(Stage) StageResult
type LoadStage func(domain.StageName) StageResult
type ExistsStage func(domain.StageName) bool

type Stage interface {
	Name() domain.StageName
	MethodContext() domain.MethodContext
	Orchestration() domain.OrchestrationName
}

type StageRequest struct {
	Name string

	Address string
	Service domain.OptionalString
	Method  domain.OptionalString

	Orchestration string
}

type StageResponse struct {
	Err domain.OptionalError
}

type SaveLink func(execute.Link) execute.LinkResult
type LoadLink func(domain.LinkName) execute.LinkResult

type SaveOrchestration func(Orchestration) OrchestrationResult
type LoadOrchestration func(domain.OrchestrationName) OrchestrationResult
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
