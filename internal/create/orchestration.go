package create

import (
	"errors"
	"github.com/DuarteMRAlves/maestro/internal"
	"github.com/DuarteMRAlves/maestro/internal/domain"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
)

type OrchestrationSaver interface {
	Save(internal.Orchestration) OrchestrationResult
}

type OrchestrationLoader interface {
	Load(internal.OrchestrationName) OrchestrationResult
}

type OrchestrationStorage interface {
	OrchestrationSaver
	OrchestrationLoader
}

type OrchestrationRequest struct {
	Name string
}

type OrchestrationResponse struct {
	Err domain.OptionalError
}

func Create(storage OrchestrationStorage) func(OrchestrationRequest) OrchestrationResponse {
	return func(req OrchestrationRequest) OrchestrationResponse {
		res := requestToOrchestration(req)
		res = BindOrchestration(verifyDupOrchestration(storage))(res)
		res = BindOrchestration(storage.Save)(res)
		return orchestrationToResponse(res)
	}
}

func requestToOrchestration(req OrchestrationRequest) OrchestrationResult {
	name, err := internal.NewOrchestrationName(req.Name)
	if err != nil {
		return ErrOrchestration(err)
	}
	if name.IsEmpty() {
		err := errdefs.InvalidArgumentWithMsg("empty orchestration name")
		return ErrOrchestration(err)
	}
	o := internal.NewOrchestration(
		name,
		[]internal.StageName{},
		[]internal.LinkName{},
	)
	return SomeOrchestration(o)
}

func updateOrchestration(
	name internal.OrchestrationName,
	loader OrchestrationLoader,
	updateFn func(internal.Orchestration) OrchestrationResult,
	saver OrchestrationSaver,
) OrchestrationResult {
	res := loader.Load(name)
	res = BindOrchestration(updateFn)(res)
	res = BindOrchestration(saver.Save)(res)
	return res
}

func verifyDupOrchestration(
	loader OrchestrationLoader,
) func(internal.Orchestration) OrchestrationResult {
	return func(o internal.Orchestration) OrchestrationResult {
		res := loader.Load(o.Name())
		if res.IsError() {
			var notFound *internal.NotFound
			err := res.Error()
			if errors.As(err, &notFound) {
				return SomeOrchestration(o)
			}
			return ErrOrchestration(err)
		}
		err := errdefs.AlreadyExistsWithMsg(
			"orchestration '%v' already exists",
			o.Name().Unwrap(),
		)
		return ErrOrchestration(err)
	}
}

func addStageNameToOrchestration(
	s internal.StageName,
) func(internal.Orchestration) internal.Orchestration {
	return func(o internal.Orchestration) internal.Orchestration {
		old := o.Stages()
		stages := make([]internal.StageName, 0, len(old)+1)
		for _, name := range old {
			stages = append(stages, name)
		}
		stages = append(stages, s)
		return internal.NewOrchestration(o.Name(), stages, o.Links())
	}
}

func addLinkNameToOrchestration(l internal.LinkName) func(internal.Orchestration) internal.Orchestration {
	return func(o internal.Orchestration) internal.Orchestration {
		old := o.Links()
		links := make([]internal.LinkName, 0, len(old)+1)
		for _, name := range old {
			links = append(links, name)
		}
		links = append(links, l)
		return internal.NewOrchestration(o.Name(), o.Stages(), links)
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
