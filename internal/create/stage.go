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

		orch, err := orchStorage.Load(orchName)
		if err != nil {
			return fmt.Errorf("add stage %s: %w", name, err)
		}

		stages := orch.Stages()
		stages = append(stages, name)
		orch = internal.NewOrchestration(orch.Name(), stages, orch.Links())

		err = orchStorage.Save(orch)
		if err != nil {
			return fmt.Errorf("add stage %s: %w", name, err)
		}
		stage := internal.NewStage(name, ctx)
		return stageStorage.Save(stage)
	}
}
