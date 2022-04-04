package create

import (
	"errors"
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal"
)

type PipelineSaver interface {
	Save(internal.Pipeline) error
}

type PipelineLoader interface {
	Load(internal.PipelineName) (internal.Pipeline, error)
}

type PipelineStorage interface {
	PipelineSaver
	PipelineLoader
}

var emptyPipelineName = errors.New("empty pipeline name")

type pipelineAlreadyExists struct{ name string }

func (err *pipelineAlreadyExists) Error() string {
	return fmt.Sprintf("pipeline '%s' already exists", err.name)
}

type unknownExecutionMode struct{ mode internal.ExecutionMode }

func (err *unknownExecutionMode) Error() string {
	return fmt.Sprintf("unknown execution mode: %s", err.mode)
}

func Pipeline(
	storage PipelineStorage,
) func(internal.PipelineName, internal.ExecutionMode) error {
	return func(name internal.PipelineName, mode internal.ExecutionMode) error {
		var modeOpt internal.PipelineOpt

		if name.IsEmpty() {
			return emptyPipelineName
		}

		_, err := storage.Load(name)
		if err == nil {
			return &pipelineAlreadyExists{name: name.Unwrap()}
		}
		var nf interface{ NotFound() }
		if !errors.As(err, &nf) {
			return err
		}

		switch mode {
		case internal.OfflineExecution:
			modeOpt = internal.WithOfflineExec()
		case internal.OnlineExecution:
			modeOpt = internal.WithOnlineExec()
		default:
			return &unknownExecutionMode{mode: mode}
		}

		p := internal.NewPipeline(name, modeOpt)
		return storage.Save(p)
	}
}
