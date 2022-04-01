package mapstore

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal"
)

type Pipelines map[internal.PipelineName]internal.Pipeline

func (s Pipelines) Save(p internal.Pipeline) error {
	s[p.Name()] = p
	return nil
}

func (s Pipelines) Load(n internal.PipelineName) (internal.Pipeline, error) {
	o, exists := s[n]
	if !exists {
		return internal.Pipeline{}, &pipelineNotFound{name: n.Unwrap()}
	}
	return o, nil
}

type pipelineNotFound struct{ name string }

func (err *pipelineNotFound) NotFound() {}

func (err *pipelineNotFound) Error() string {
	return fmt.Sprintf("pipeline not found: %s", err.name)
}
