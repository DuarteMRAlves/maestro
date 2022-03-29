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
	emptyStageName = errors.New("empty stage name")
	emptyAddress   = errors.New("empty address")
)

type stageAlreadyExists struct{ name string }

func (err *stageAlreadyExists) Error() string {
	return fmt.Sprintf("stage '%s' already exists", err.name)
}

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
			return emptyStageName
		}
		addr := ctx.Address()
		if addr.IsEmpty() {
			return emptyAddress
		}

		if orchName.IsEmpty() {
			return emptyOrchestrationName
		}

		_, err := stageStorage.Load(name)
		if err == nil {
			return &stageAlreadyExists{name: name.Unwrap()}
		}
		var nf interface{ NotFound() }
		if !errors.As(err, &nf) {
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
