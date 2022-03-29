package mapstore

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal"
)

type Stages map[internal.StageName]internal.Stage

func (s Stages) Save(o internal.Stage) error {
	s[o.Name()] = o
	return nil
}

func (s Stages) Load(n internal.StageName) (internal.Stage, error) {
	o, exists := s[n]
	if !exists {
		return internal.Stage{}, &stageNotFound{name: n.Unwrap()}
	}
	return o, nil
}

type stageNotFound struct{ name string }

func (err *stageNotFound) NotFound() {}

func (err *stageNotFound) Error() string {
	return fmt.Sprintf("stage not found: %s", err.name)
}
