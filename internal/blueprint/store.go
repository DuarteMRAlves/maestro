package blueprint

import (
	"errors"
	"sync"
)

type Store interface {
	Create(b *Blueprint) error
}

type store struct {
	blueprints sync.Map
}

func NewStore() Store {
	return &store{blueprints: sync.Map{}}
}

func (st *store) Create(config *Blueprint) error {
	if ok, err := verifyConfig(config); !ok {
		return err
	}

	bp := config.Clone()
	_, prev := st.blueprints.LoadOrStore(bp.Name, bp)
	if prev {
		return AlreadyExists{Name: bp.Name}
	}
	return nil
}

func verifyConfig(config *Blueprint) (bool, error) {
	if config == nil {
		return false, errors.New("nil config")
	}
	if len(config.Stages) != 0 {
		return false, errors.New("blueprint should not have Stages")
	}
	if len(config.Links) != 0 {
		return false, errors.New("blueprint should not have Links")
	}
	return true, nil
}
