package blueprint

import (
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/validate"
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
	if ok, err := validate.ArgNotNil(config, "config"); !ok {
		return err
	}

	bp := config.Clone()
	_, prev := st.blueprints.LoadOrStore(bp.Name, bp)
	if prev {
		return errdefs.AlreadyExistsWithMsg(
			"blueprint '%v' already exists",
			bp.Name)
	}
	return nil
}
