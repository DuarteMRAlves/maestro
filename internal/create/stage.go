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
	stageStorage StageStorage, pipelineStorage PipelineStorage,
) func(internal.StageName, internal.MethodContext, internal.PipelineName) error {
	return func(
		name internal.StageName,
		methodContext internal.MethodContext,
		pipelineName internal.PipelineName,
	) error {
		if name.IsEmpty() {
			return emptyStageName
		}
		addr := methodContext.Address()
		if addr.IsEmpty() {
			return emptyAddress
		}

		if pipelineName.IsEmpty() {
			return emptyPipelineName
		}

		_, err := stageStorage.Load(name)
		if err == nil {
			return &stageAlreadyExists{name: name.Unwrap()}
		}
		var nf interface{ NotFound() }
		if !errors.As(err, &nf) {
			return err
		}

		pipeline, err := pipelineStorage.Load(pipelineName)
		if err != nil {
			return fmt.Errorf("add stage %s: %w", name, err)
		}

		stages := pipeline.Stages()
		stages = append(stages, name)
		pipeline = internal.NewPipeline(pipeline.Name(), stages, pipeline.Links())

		err = pipelineStorage.Save(pipeline)
		if err != nil {
			return fmt.Errorf("add stage %s: %w", name, err)
		}
		stage := internal.NewStage(name, methodContext)
		return stageStorage.Save(stage)
	}
}
