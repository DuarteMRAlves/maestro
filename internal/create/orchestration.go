package create

import (
	"github.com/DuarteMRAlves/maestro/internal/domain"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
)

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

func CreateOrchestration(
	existsFn ExistsOrchestration,
	saveFn SaveOrchestration,
) func(OrchestrationRequest) OrchestrationResponse {
	return func(req OrchestrationRequest) OrchestrationResponse {
		res := requestToOrchestration(req)
		res = BindOrchestration(verifyDupOrchestration(existsFn))(res)
		res = BindOrchestration(saveFn)(res)
		return orchestrationToResponse(res)
	}
}

func requestToOrchestration(req OrchestrationRequest) OrchestrationResult {
	name, err := domain.NewOrchestrationName(req.Name)
	if err != nil {
		return ErrOrchestration(err)
	}
	o := NewOrchestration(name, []domain.StageName{}, []domain.LinkName{})
	return SomeOrchestration(o)
}

func verifyDupOrchestration(
	existsFn ExistsOrchestration,
) func(Orchestration) OrchestrationResult {
	return func(o Orchestration) OrchestrationResult {
		if existsFn(o.Name()) {
			err := errdefs.AlreadyExistsWithMsg(
				"orchestration '%v' already exists",
				o.Name().Unwrap(),
			)
			return ErrOrchestration(err)
		}
		return SomeOrchestration(o)
	}
}

func orchestrationToResponse(res OrchestrationResult) OrchestrationResponse {
	var errOpt domain.OptionalError
	if res.IsError() {
		errOpt = domain.NewPresentError(res.Error())
	} else {
		errOpt = domain.NewEmptyError()
	}
	return OrchestrationResponse{Err: errOpt}
}
