package create

import (
	"github.com/DuarteMRAlves/maestro/internal/domain"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
)

func CreateOrchestration(
	existsFn ExistsOrchestration,
	saveFn SaveOrchestration,
) func(OrchestrationRequest) OrchestrationResponse {
	return func(req OrchestrationRequest) OrchestrationResponse {
		res := requestToOrchestration(req)
		res = domain.BindOrchestration(verifyDupOrchestration(existsFn))(res)
		res = domain.BindOrchestration(saveFn)(res)
		return orchestrationToResponse(res)
	}
}

func requestToOrchestration(req OrchestrationRequest) domain.OrchestrationResult {
	name, err := domain.NewOrchestrationName(req.Name)
	if err != nil {
		return domain.ErrOrchestration(err)
	}
	o := domain.NewOrchestration(name, []domain.Stage{}, []domain.Link{})
	return domain.SomeOrchestration(o)
}

func verifyDupOrchestration(
	existsFn ExistsOrchestration,
) func(domain.Orchestration) domain.OrchestrationResult {
	return func(o domain.Orchestration) domain.OrchestrationResult {
		if existsFn(o.Name()) {
			err := errdefs.AlreadyExistsWithMsg(
				"orchestration '%v' already exists",
				o.Name().Unwrap(),
			)
			return domain.ErrOrchestration(err)
		}
		return domain.SomeOrchestration(o)
	}
}

func orchestrationToResponse(res domain.OrchestrationResult) OrchestrationResponse {
	var errOpt domain.OptionalError
	if res.IsError() {
		errOpt = domain.NewPresentError(res.Error())
	} else {
		errOpt = domain.NewEmptyError()
	}
	return OrchestrationResponse{Err: errOpt}
}
