package create

import "github.com/DuarteMRAlves/maestro/internal/domain"

type OrchestrationSaver interface {
	Save(Orchestration) OrchestrationResult
}

type OrchestrationLoader interface {
	Load(domain.OrchestrationName) OrchestrationResult
}

type OrchestrationExistsVerifier interface {
	Verify(domain.OrchestrationName) bool
}

type OrchestrationStorage interface {
	OrchestrationSaver
	OrchestrationLoader
	OrchestrationExistsVerifier
}

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

type orchestration struct {
	name   domain.OrchestrationName
	stages []domain.StageName
	links  []domain.LinkName
}

func (o orchestration) Name() domain.OrchestrationName {
	return o.name
}

func (o orchestration) Stages() []domain.StageName {
	return o.stages
}

func (o orchestration) Links() []domain.LinkName {
	return o.links
}

func NewOrchestration(
	name domain.OrchestrationName,
	stages []domain.StageName,
	links []domain.LinkName,
) Orchestration {
	return &orchestration{
		name:   name,
		stages: stages,
		links:  links,
	}
}

func updateOrchestration(
	name domain.OrchestrationName,
	loader OrchestrationLoader,
	updateFn func(Orchestration) OrchestrationResult,
	saver OrchestrationSaver,
) OrchestrationResult {
	res := loader.Load(name)
	res = BindOrchestration(updateFn)(res)
	res = BindOrchestration(saver.Save)(res)
	return res
}

func addStageNameToOrchestration(
	s domain.StageName,
) func(Orchestration) Orchestration {
	return func(o Orchestration) Orchestration {
		old := o.Stages()
		stages := make([]domain.StageName, 0, len(old)+1)
		for _, name := range old {
			stages = append(stages, name)
		}
		stages = append(stages, s)
		return NewOrchestration(o.Name(), stages, o.Links())
	}
}

func addLinkNameToOrchestration(l domain.LinkName) func(Orchestration) Orchestration {
	return func(o Orchestration) Orchestration {
		old := o.Links()
		links := make([]domain.LinkName, 0, len(old)+1)
		for _, name := range old {
			links = append(links, name)
		}
		links = append(links, l)
		return NewOrchestration(o.Name(), o.Stages(), links)
	}
}
