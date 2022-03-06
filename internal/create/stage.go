package create

import (
	"github.com/DuarteMRAlves/maestro/internal/domain"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
)

type stage struct {
	name          domain.StageName
	methodCtx     domain.MethodContext
	orchestration domain.OrchestrationName
}

func (s stage) Name() domain.StageName {
	return s.name
}

func (s stage) MethodContext() domain.MethodContext {
	return s.methodCtx
}

func (s stage) Orchestration() domain.OrchestrationName {
	return s.orchestration
}

func NewStage(
	name domain.StageName,
	methodCtx domain.MethodContext,
	orchestration domain.OrchestrationName,
) Stage {
	return stage{
		name:          name,
		methodCtx:     methodCtx,
		orchestration: orchestration,
	}
}

func CreateStage(
	existsStage ExistsStage,
	saveStage SaveStage,
	existsOrchestration ExistsOrchestration,
	loadOrchestration LoadOrchestration,
	saveOrchestration SaveOrchestration,
) func(StageRequest) StageResponse {
	return func(req StageRequest) StageResponse {
		res := requestToStage(req)
		res = BindStage(verifyDupStage(existsStage))(res)
		res = BindStage(verifyExistsOrchestration(existsOrchestration))(res)
		res = BindStage(addStage(loadOrchestration, saveOrchestration))(res)
		res = BindStage(saveStage)(res)
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
	addr, err := domain.NewAddress(req.Address)
	if err != nil {
		return ErrStage(err)
	}

	if req.Service.Present() {
		service, err := domain.NewService(req.Service.Unwrap())
		if err != nil {
			return ErrStage(err)
		}
		serviceOpt = domain.NewPresentService(service)
	}

	if req.Method.Present() {
		method, err := domain.NewMethod(req.Method.Unwrap())
		if err != nil {
			return ErrStage(err)
		}
		methodOpt = domain.NewPresentMethod(method)
	}

	ctx := domain.NewMethodContext(addr, serviceOpt, methodOpt)

	orchestrationName, err := domain.NewOrchestrationName(req.Orchestration)
	if err != nil {
		return ErrStage(err)
	}

	return SomeStage(NewStage(name, ctx, orchestrationName))
}

func verifyDupStage(existsFn ExistsStage) func(Stage) StageResult {
	return func(s Stage) StageResult {
		if existsFn(s.Name()) {
			err := errdefs.AlreadyExistsWithMsg(
				"stage '%v' already exists",
				s.Name().Unwrap(),
			)
			return ErrStage(err)
		}
		return SomeStage(s)
	}
}

func verifyExistsOrchestration(existsFn ExistsOrchestration) func(Stage) StageResult {
	return func(s Stage) StageResult {
		if !existsFn(s.Orchestration()) {
			err := errdefs.NotFoundWithMsg(
				"orchestration '%v' not found",
				s.Orchestration().Unwrap(),
			)
			return ErrStage(err)
		}
		return SomeStage(s)
	}
}

func addStage(
	loadFn LoadOrchestration,
	saveFn SaveOrchestration,
) func(Stage) StageResult {
	return func(stage Stage) StageResult {
		name := stage.Orchestration()
		updateFn := addStageToOrchestration(stage)
		res := updateOrchestration(name, loadFn, updateFn, saveFn)
		if res.IsError() {
			err := errdefs.PrependMsg(
				res.Error(),
				"add stage %s to orchestration: %s",
				stage.Name(),
				stage.Orchestration(),
			)
			return ErrStage(err)
		}
		return SomeStage(stage)
	}
}

func updateOrchestration(
	name domain.OrchestrationName,
	loadFn LoadOrchestration,
	updateFn func(Orchestration) OrchestrationResult,
	saveFn SaveOrchestration,
) OrchestrationResult {
	res := loadFn(name)
	res = BindOrchestration(updateFn)(res)
	res = BindOrchestration(saveFn)(res)
	return res
}

func addStageToOrchestration(stage Stage) func(Orchestration) OrchestrationResult {
	return func(o Orchestration) OrchestrationResult {
		old := o.Stages()
		stages := make([]domain.StageName, 0, len(old)+1)
		for _, s := range old {
			stages = append(stages, s)
		}
		stages = append(stages, stage.Name())
		updated := NewOrchestration(o.Name(), stages, o.Links())
		return SomeOrchestration(updated)
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
