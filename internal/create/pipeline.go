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

func Pipeline(storage PipelineStorage) func(internal.PipelineName) error {
	return func(name internal.PipelineName) error {
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

		p := internal.NewPipeline(name, nil, nil)
		return storage.Save(p)
	}
}
