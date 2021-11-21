package stage

import (
	"github.com/DuarteMRAlves/maestro/internal/assert"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"sync"
)

type Store interface {
	Create(s *Stage) error
}

type store struct {
	stages sync.Map
}

func NewStore() Store {
	return &store{stages: sync.Map{}}
}

func (st *store) Create(config *Stage) error {
	if ok, err := assert.ArgNotNil(config, "config"); !ok {
		return err
	}

	s := config.Clone()
	_, prev := st.stages.LoadOrStore(s.Name, s)
	if prev {
		return errdefs.AlreadyExistsWithMsg("stage '%v' already exists", s.Name)
	}
	return nil
}
