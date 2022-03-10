package create

import "github.com/DuarteMRAlves/maestro/internal/domain"

type StageSaver interface {
	Save(Stage) StageResult
}

type StageLoader interface {
	Load(domain.StageName) StageResult
}

type StageStorage interface {
	StageSaver
	StageLoader
}

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

type stage struct {
	name          domain.StageName
	methodCtx     domain.MethodContext
	orchestration domain.OrchestrationName
}

func (s stage) Name() domain.StageName {
	return s.name
}

func (s stage) MethodContext() domain.MethodContext {
	return s.methodCtx
}

func (s stage) Orchestration() domain.OrchestrationName {
	return s.orchestration
}

func NewStage(
	name domain.StageName,
	methodCtx domain.MethodContext,
	orchestration domain.OrchestrationName,
) Stage {
	return stage{
		name:          name,
		methodCtx:     methodCtx,
		orchestration: orchestration,
	}
}
