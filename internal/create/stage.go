package create

import (
	"errors"
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal"
	"github.com/DuarteMRAlves/maestro/internal/domain"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
)

type StageSaver interface {
	Save(internal.Stage) error
}

type StageLoader interface {
	Load(internal.StageName) (internal.Stage, error)
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
		serviceOpt := internal.NewEmptyService()
		methodOpt := internal.NewEmptyMethod()

		name, err := internal.NewStageName(req.Name)
		if err != nil {
			return StageResponse{Err: domain.NewPresentError(err)}
		}
		if name.IsEmpty() {
			return StageResponse{Err: domain.NewPresentError(EmptyStageName)}
		}
		addr := internal.NewAddress(req.Address)
		if addr.IsEmpty() {
			return StageResponse{Err: domain.NewPresentError(EmptyAddress)}
		}

		if req.Service.Present() {
			service := internal.NewService(req.Service.Unwrap())
			if service.IsEmpty() {
				return StageResponse{Err: domain.NewPresentError(EmptyService)}
			}
			serviceOpt = internal.NewPresentService(service)
		}

		if req.Method.Present() {
			method := internal.NewMethod(req.Method.Unwrap())
			if method.IsEmpty() {
				return StageResponse{Err: domain.NewPresentError(EmptyMethod)}
			}
			methodOpt = internal.NewPresentMethod(method)
		}

		ctx := internal.NewMethodContext(addr, serviceOpt, methodOpt)

		orchestrationName, err := internal.NewOrchestrationName(req.Orchestration)
		if err != nil {
			return StageResponse{Err: domain.NewPresentError(err)}
		}

		_, err = stageStorage.Load(name)
		if err == nil {
			err := errdefs.AlreadyExistsWithMsg(
				"stage '%v' already exists",
				name.Unwrap(),
			)
			return StageResponse{Err: domain.NewPresentError(err)}
		}
		var notFound *internal.NotFound
		if !errors.As(err, &notFound) {
			return StageResponse{Err: domain.NewPresentError(err)}
		}
		_, err = orchStorage.Load(orchestrationName)
		if err != nil {
			return StageResponse{Err: domain.NewPresentError(err)}
		}
		updateFn := addStageNameToOrchestration(name)
		err = updateOrchestration(
			orchestrationName,
			orchStorage,
			updateFn,
			orchStorage,
		)
		if err != nil {
			err := errdefs.PrependMsg(
				err,
				"add stage %s to orchestration: %s",
				name,
				orchestrationName,
			)
			return StageResponse{Err: domain.NewPresentError(err)}
		}
		stage := internal.NewStage(name, ctx, orchestrationName)
		err = stageStorage.Save(stage)
		errOpt := domain.NewEmptyError()
		if err != nil {
			errOpt = domain.NewPresentError(err)
		}
		return StageResponse{Err: errOpt}
	}
}
