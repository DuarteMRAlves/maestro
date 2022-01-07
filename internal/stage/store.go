package stage

import (
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/validate"
	"sync"
)

type Store interface {
	Create(s *Stage) error
	Contains(name string) bool
	GetByName(name string) (*Stage, bool)
	GetMatching(query *apitypes.Stage) []*Stage
}

type store struct {
	stages sync.Map
}

func NewStore() Store {
	return &store{stages: sync.Map{}}
}

func (st *store) Create(config *Stage) error {
	if ok, err := validate.ArgNotNil(config, "config"); !ok {
		return err
	}

	s := config.Clone()
	_, prev := st.stages.LoadOrStore(s.name, s)
	if prev {
		return errdefs.AlreadyExistsWithMsg("stage '%v' already exists", s.name)
	}
	return nil
}

func (st *store) Contains(name string) bool {
	_, ok := st.stages.Load(name)
	return ok
}

func (st *store) GetByName(name string) (*Stage, bool) {
	loaded, ok := st.stages.Load(name)
	if !ok {
		return nil, false
	}
	stage, ok := loaded.(*Stage)
	return stage, ok
}

func (st *store) GetMatching(query *apitypes.Stage) []*Stage {
	if query == nil {
		query = &apitypes.Stage{}
	}
	filter := buildQueryFilter(query)
	res := make([]*Stage, 0)
	st.stages.Range(
		func(key, value interface{}) bool {
			s, ok := value.(*Stage)
			if !ok {
				return false
			}
			if filter(s) {
				res = append(res, s.Clone())
			}
			return true
		})
	return res
}

func buildQueryFilter(query *apitypes.Stage) func(s *Stage) bool {
	filters := make([]func(s *Stage) bool, 0)
	if query.Name != "" {
		filters = append(
			filters,
			func(s *Stage) bool {
				return s.name == query.Name
			})
	}
	if query.Phase != "" {
		filters = append(
			filters,
			func(s *Stage) bool {
				return s.phase == query.Phase
			})
	}
	if query.Asset != "" {
		filters = append(
			filters,
			func(s *Stage) bool {
				return s.asset == query.Asset
			})
	}
	if query.Service != "" {
		filters = append(
			filters,
			func(s *Stage) bool {
				return s.rpc.Service().Name() == query.Service
			})
	}
	if query.Rpc != "" {
		filters = append(
			filters,
			func(s *Stage) bool {
				return s.rpc.Name() == query.Rpc
			})
	}
	if query.Address != "" {
		filters = append(
			filters,
			func(s *Stage) bool {
				return s.address == query.Address
			})
	}
	switch len(filters) {
	case 0:
		return func(s *Stage) bool { return true }
	case 1:
		return filters[0]
	default:
		return func(s *Stage) bool {
			for _, f := range filters {
				if !f(s) {
					return false
				}
			}
			return true
		}
	}
}
