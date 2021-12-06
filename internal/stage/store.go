package stage

import (
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/validate"
	"sync"
)

type Store interface {
	Create(s *Stage) error
	Contains(name string) bool
	Get(query *Stage) []*Stage
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
	_, prev := st.stages.LoadOrStore(s.Name, s)
	if prev {
		return errdefs.AlreadyExistsWithMsg("stage '%v' already exists", s.Name)
	}
	return nil
}

func (st *store) Contains(name string) bool {
	_, ok := st.stages.Load(name)
	return ok
}

func (st *store) Get(query *Stage) []*Stage {
	if query == nil {
		query = &Stage{}
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

func buildQueryFilter(query *Stage) func(s *Stage) bool {
	filters := make([]func(s *Stage) bool, 0)
	if query.Name != "" {
		filters = append(
			filters,
			func(s *Stage) bool {
				return s.Name == query.Name
			})
	}
	if query.Asset != "" {
		filters = append(
			filters,
			func(s *Stage) bool {
				return s.Asset == query.Asset
			})
	}
	if query.Service != "" {
		filters = append(
			filters,
			func(s *Stage) bool {
				return s.Service == query.Service
			})
	}
	if query.Method != "" {
		filters = append(
			filters,
			func(s *Stage) bool {
				return s.Method == query.Method
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
