package execute

import "github.com/DuarteMRAlves/maestro/internal/domain"

type LoadLink func(domain.LinkName) LinkResult
type LoadOrchestration func(domain.OrchestrationName) OrchestrationResult

type LinkEndpoint interface {
	Stage() domain.Stage
	Field() domain.OptionalMessageField
}

type Link interface {
	Name() domain.LinkName
	Source() LinkEndpoint
	Target() LinkEndpoint
}

type Orchestration interface {
	Name() domain.OrchestrationName
	Stages() []domain.Stage
	Links() []Link
}
