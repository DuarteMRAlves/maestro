package blueprint

import (
	"errors"
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/identifier"
	"sync"
)

type Store interface {
	Create(b *Blueprint) (identifier.Id, error)
}

type store struct {
	blueprints sync.Map
	gen        identifier.Generator
}

func NewStore() Store {
	return &store{
		blueprints: sync.Map{},
		gen:        identifier.GenForSize(IdSize),
	}
}

func (st *store) Create(config *Blueprint) (identifier.Id, error) {
	if ok, err := verifyConfig(config); !ok {
		return identifier.Empty(), err
	}
	bp := config.Clone()

	id, err := st.gen()
	if err != nil {
		return identifier.Empty(), fmt.Errorf("store create blueprint: %v", err)
	}
	st.blueprints.Store(id, bp)

	bp.Id = id

	return id.Clone(), nil
}

func verifyConfig(config *Blueprint) (bool, error) {
	if config == nil {
		return false, errors.New("nil config")
	}
	if !config.Id.IsEmpty() {
		return false, errors.New("blueprint identifier should not be defined")
	}
	if len(config.stages) != 0 {
		return false, errors.New("blueprint should not have assets")
	}
	if len(config.links) != 0 {
		return false, errors.New("blueprint should not have links")
	}
	return true, nil
}
