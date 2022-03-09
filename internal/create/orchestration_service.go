package create

import (
	"github.com/DuarteMRAlves/maestro/internal/domain"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
)

func CreateOrchestration(storage OrchestrationStorage) func(OrchestrationRequest) OrchestrationResponse {
	return func(req OrchestrationRequest) OrchestrationResponse {
		res := requestToOrchestration(req)
		res = BindOrchestration(verifyDupOrchestration(storage))(res)
		res = BindOrchestration(storage.Save)(res)
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
	verifier OrchestrationExistsVerifier,
) func(Orchestration) OrchestrationResult {
	return func(o Orchestration) OrchestrationResult {
		if verifier.Verify(o.Name()) {
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
