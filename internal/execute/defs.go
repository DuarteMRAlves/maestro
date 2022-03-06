package execute

import "github.com/DuarteMRAlves/maestro/internal/domain"

type LoadLink func(domain.LinkName) LinkResult
type LoadOrchestration func(domain.OrchestrationName) OrchestrationResult

type Stage interface {
	Name() domain.StageName
	MethodContext() domain.MethodContext
}

type LinkEndpoint interface {
	Stage() Stage
	Field() domain.OptionalMessageField
}

type Link interface {
	Name() domain.LinkName
	Source() LinkEndpoint
	Target() LinkEndpoint
}

type Orchestration interface {
	Name() domain.OrchestrationName
	Stages() []Stage
	Links() []Link
}
