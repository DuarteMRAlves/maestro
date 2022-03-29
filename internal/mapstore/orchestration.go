package mapstore

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal"
)

type Orchestrations map[internal.OrchestrationName]internal.Orchestration

func (s Orchestrations) Save(o internal.Orchestration) error {
	s[o.Name()] = o
	return nil
}

func (s Orchestrations) Load(
	n internal.OrchestrationName,
) (internal.Orchestration, error) {
	o, exists := s[n]
	if !exists {
		return internal.Orchestration{}, &orchestrationNotFound{name: n.Unwrap()}
	}
	return o, nil
}

type orchestrationNotFound struct{ name string }

func (err *orchestrationNotFound) NotFound() {}

func (err *orchestrationNotFound) Error() string {
	return fmt.Sprintf("orchestration not found: %s", err.name)
}
