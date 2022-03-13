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

var (
	EmptyStageName = fmt.Errorf("empty stage name")
	EmptyAddress   = fmt.Errorf("empty address")
	EmptyService   = fmt.Errorf("empty service")
	EmptyMethod    = fmt.Errorf("empty method")
)

func Stage(
	stageStorage StageStorage,
	orchStorage OrchestrationStorage,
) func(StageRequest) error {
	return func(req StageRequest) error {
		serviceOpt := internal.NewEmptyService()
		methodOpt := internal.NewEmptyMethod()

		name, err := internal.NewStageName(req.Name)
		if err != nil {
			return err
		}
		if name.IsEmpty() {
			return EmptyStageName
		}
		addr := internal.NewAddress(req.Address)
		if addr.IsEmpty() {
			return EmptyAddress
		}

		if req.Service.Present() {
			service := internal.NewService(req.Service.Unwrap())
			if service.IsEmpty() {
				return EmptyService
			}
			serviceOpt = internal.NewPresentService(service)
		}

		if req.Method.Present() {
			method := internal.NewMethod(req.Method.Unwrap())
			if method.IsEmpty() {
				return EmptyMethod
			}
			methodOpt = internal.NewPresentMethod(method)
		}

		ctx := internal.NewMethodContext(addr, serviceOpt, methodOpt)

		orchestrationName, err := internal.NewOrchestrationName(req.Orchestration)
		if err != nil {
			return err
		}

		_, err = stageStorage.Load(name)
		if err == nil {
			return &internal.AlreadyExists{Type: "stage", Ident: name.Unwrap()}
		}
		var notFound *internal.NotFound
		if !errors.As(err, &notFound) {
			return err
		}
		_, err = orchStorage.Load(orchestrationName)
		if err != nil {
			return err
		}
		updateFn := addStageNameToOrchestration(name)
		err = updateOrchestration(
			orchestrationName,
			orchStorage,
			updateFn,
			orchStorage,
		)
		if err != nil {
			return errdefs.PrependMsg(
				err,
				"add stage %s to orchestration: %s",
				name,
				orchestrationName,
			)
		}
		stage := internal.NewStage(name, ctx, orchestrationName)
		return stageStorage.Save(stage)
	}
}
