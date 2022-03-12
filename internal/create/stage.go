package create

import (
	"github.com/DuarteMRAlves/maestro/internal/domain"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
)

type StageSaver interface {
	Save(domain.Stage) StageResult
}

type StageLoader interface {
	Load(domain.StageName) StageResult
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

func CreateStage(
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
	serviceOpt := domain.NewEmptyService()
	methodOpt := domain.NewEmptyMethod()

	name, err := domain.NewStageName(req.Name)
	if err != nil {
		return ErrStage(err)
	}
	if name.IsEmpty() {
		err := errdefs.InvalidArgumentWithMsg("empty stage name")
		return ErrStage(err)
	}
	addr := domain.NewAddress(req.Address)
	if addr.IsEmpty() {
		return ErrStage(errdefs.InvalidArgumentWithMsg("empty address"))
	}

	if req.Service.Present() {
		service := domain.NewService(req.Service.Unwrap())
		if service.IsEmpty() {
			err := errdefs.InvalidArgumentWithMsg("empty service")
			return ErrStage(err)
		}
		serviceOpt = domain.NewPresentService(service)
	}

	if req.Method.Present() {
		method := domain.NewMethod(req.Method.Unwrap())
		if method.IsEmpty() {
			err := errdefs.InvalidArgumentWithMsg("empty method")
			return ErrStage(err)
		}
		methodOpt = domain.NewPresentMethod(method)
	}

	ctx := domain.NewMethodContext(addr, serviceOpt, methodOpt)

	orchestrationName, err := domain.NewOrchestrationName(req.Orchestration)
	if err != nil {
		return ErrStage(err)
	}

	return SomeStage(domain.NewStage(name, ctx, orchestrationName))
}

func verifyDupStage(loader StageLoader) func(domain.Stage) StageResult {
	return func(s domain.Stage) StageResult {
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

func verifyExistsOrchestration(orchLoader OrchestrationLoader) func(domain.Stage) StageResult {
	return func(s domain.Stage) StageResult {
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
) func(domain.Stage) StageResult {
	return func(s domain.Stage) StageResult {
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
