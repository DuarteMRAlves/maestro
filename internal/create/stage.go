package create

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal"
	"github.com/DuarteMRAlves/maestro/internal/domain"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
)

type StageSaver interface {
	Save(internal.Stage) StageResult
}

type StageLoader interface {
	Load(internal.StageName) StageResult
}

type StageStorage interface {
	StageSaver
	StageLoader
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

var (
	EmptyStageName = fmt.Errorf("empty stage name")
	EmptyAddress   = fmt.Errorf("empty address")
	EmptyService   = fmt.Errorf("empty service")
	EmptyMethod    = fmt.Errorf("empty method")
)

func Stage(
	stageStorage StageStorage,
	orchStorage OrchestrationStorage,
) func(StageRequest) StageResponse {
	return func(req StageRequest) StageResponse {
		res := requestToStage(req)
		res = BindStage(verifyDupStage(stageStorage))(res)
		res = BindStage(verifyExistsOrchestration(orchStorage))(res)
		res = BindStage(addStage(orchStorage, orchStorage))(res)
		res = BindStage(stageStorage.Save)(res)
		return stageToResponse(res)
	}
}

func requestToStage(req StageRequest) StageResult {
	serviceOpt := internal.NewEmptyService()
	methodOpt := internal.NewEmptyMethod()

	name, err := internal.NewStageName(req.Name)
	if err != nil {
		return ErrStage(err)
	}
	if name.IsEmpty() {
		return ErrStage(EmptyStageName)
	}
	addr := internal.NewAddress(req.Address)
	if addr.IsEmpty() {
		return ErrStage(EmptyAddress)
	}

	if req.Service.Present() {
		service := internal.NewService(req.Service.Unwrap())
		if service.IsEmpty() {
			return ErrStage(EmptyService)
		}
		serviceOpt = internal.NewPresentService(service)
	}

	if req.Method.Present() {
		method := internal.NewMethod(req.Method.Unwrap())
		if method.IsEmpty() {
			return ErrStage(EmptyMethod)
		}
		methodOpt = internal.NewPresentMethod(method)
	}

	ctx := internal.NewMethodContext(addr, serviceOpt, methodOpt)

	orchestrationName, err := domain.NewOrchestrationName(req.Orchestration)
	if err != nil {
		return ErrStage(err)
	}

	return SomeStage(internal.NewStage(name, ctx, orchestrationName))
}

func verifyDupStage(loader StageLoader) func(internal.Stage) StageResult {
	return func(s internal.Stage) StageResult {
		res := loader.Load(s.Name())
		if res.IsError() {
			err := res.Error()
			if errdefs.IsNotFound(err) {
				return SomeStage(s)
			}
			return ErrStage(err)
		}
		err := errdefs.AlreadyExistsWithMsg(
			"stage '%v' already exists",
			s.Name().Unwrap(),
		)
		return ErrStage(err)
	}
}

func verifyExistsOrchestration(orchLoader OrchestrationLoader) func(internal.Stage) StageResult {
	return func(s internal.Stage) StageResult {
		res := orchLoader.Load(s.Orchestration())
		if res.IsError() {
			return ErrStage(res.Error())
		}
		return SomeStage(s)
	}
}

func addStage(
	loader OrchestrationLoader,
	saver OrchestrationSaver,
) func(internal.Stage) StageResult {
	return func(s internal.Stage) StageResult {
		name := s.Orchestration()
		updateFn := ReturnOrchestration(addStageNameToOrchestration(s.Name()))
		res := updateOrchestration(name, loader, updateFn, saver)
		if res.IsError() {
			err := errdefs.PrependMsg(
				res.Error(),
				"add stage %s to orchestration: %s",
				s.Name(),
				s.Orchestration(),
			)
			return ErrStage(err)
		}
		return SomeStage(s)
	}
}

func stageToResponse(res StageResult) StageResponse {
	var errOpt domain.OptionalError
	if res.IsError() {
		errOpt = domain.NewPresentError(res.Error())
	} else {
		errOpt = domain.NewEmptyError()
	}
	return StageResponse{Err: errOpt}
}
