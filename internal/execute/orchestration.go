package execute

import (
	"github.com/DuarteMRAlves/maestro/internal/domain"
)

type orchestration struct {
	name   domain.OrchestrationName
	stages []domain.Stage
	links  []Link
}

func (o orchestration) Name() domain.OrchestrationName {
	return o.name
}

func (o orchestration) Stages() []domain.Stage {
	return o.stages
}

func (o orchestration) Links() []Link {
	return o.links
}

func NewOrchestration(
	name domain.OrchestrationName,
	stages []domain.Stage,
	links []Link,
) Orchestration {
	return &orchestration{
		name:   name,
		stages: stages,
		links:  links,
	}
}
