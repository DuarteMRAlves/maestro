package create

import (
	"errors"
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal"
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

var (
	EmptyStageName = fmt.Errorf("empty stage name")
	EmptyAddress   = fmt.Errorf("empty address")
	EmptyService   = fmt.Errorf("empty service")
	EmptyMethod    = fmt.Errorf("empty method")
)

func Stage(
	stageStorage StageStorage,
	orchStorage OrchestrationStorage,
) func(
	internal.StageName,
	internal.MethodContext,
	internal.OrchestrationName,
) error {
	return func(
		name internal.StageName,
		ctx internal.MethodContext,
		orchName internal.OrchestrationName,
	) error {
		if name.IsEmpty() {
			return EmptyStageName
		}
		addr := ctx.Address()
		if addr.IsEmpty() {
			return EmptyAddress
		}

		service := ctx.Service()
		if service.Present() && service.Unwrap().IsEmpty() {
			return EmptyService
		}

		method := ctx.Method()
		if method.Present() && method.Unwrap().IsEmpty() {
			return EmptyMethod
		}

		if orchName.IsEmpty() {
			return EmptyOrchestrationName
		}

		_, err := stageStorage.Load(name)
		if err == nil {
			return &internal.AlreadyExists{Type: "stage", Ident: name.Unwrap()}
		}
		var notFound *internal.NotFound
		if !errors.As(err, &notFound) {
			return err
		}

		_, err = orchStorage.Load(orchName)
		if err != nil {
			return err
		}
		updateFn := addStageNameToOrchestration(name)
		err = updateOrchestration(orchName, orchStorage, updateFn, orchStorage)
		if err != nil {
			format := "add stage %s to orchestration %s: %w"
			return fmt.Errorf(format, name, orchName, err)
		}
		stage := internal.NewStage(name, ctx, orchName)
		return stageStorage.Save(stage)
	}
}
