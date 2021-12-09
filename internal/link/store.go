package link

import (
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/validate"
	"sync"
)

type Store interface {
	Create(l *Link) error
	Contains(name string) bool
	Get(query *Link) []*Link
}

type store struct {
	links sync.Map
}

func NewStore() Store {
	return &store{links: sync.Map{}}
}

func (st *store) Create(config *Link) error {
	if ok, err := validate.ArgNotNil(config, "config"); !ok {
		return err
	}

	l := config.Clone()
	_, prev := st.links.LoadOrStore(l.Name, l)
	if prev {
		return errdefs.AlreadyExistsWithMsg("link '%v' already exists", l.Name)
	}
	return nil
}

// Contains returns true if a link with the given name exists and false
// otherwise.
func (st *store) Contains(name string) bool {
	_, ok := st.links.Load(name)
	return ok
}

func (st *store) Get(query *Link) []*Link {
	if query == nil {
		query = &Link{}
	}
	filter := buildQueryFilter(query)
	res := make([]*Link, 0)
	st.links.Range(
		func(key, value interface{}) bool {
			l, ok := value.(*Link)
			if !ok {
				return false
			}
			if filter(l) {
				res = append(res, l.Clone())
			}
			return true
		})
	return res
}

func buildQueryFilter(query *Link) func(l *Link) bool {
	filters := make([]func(l *Link) bool, 0)
	if query.Name != "" {
		filters = append(
			filters,
			func(l *Link) bool {
				return l.Name == query.Name
			})
	}
	if query.SourceStage != "" {
		filters = append(
			filters,
			func(l *Link) bool {
				return l.SourceStage == query.SourceStage
			})
	}
	if query.SourceField != "" {
		filters = append(
			filters,
			func(l *Link) bool {
				return l.SourceField == query.SourceField
			})
	}
	if query.TargetStage != "" {
		filters = append(
			filters,
			func(l *Link) bool {
				return l.TargetStage == query.TargetStage
			})
	}
	switch len(filters) {
	case 0:
		return func(l *Link) bool { return true }
	case 1:
		return filters[0]
	default:
		return func(l *Link) bool {
			for _, f := range filters {
				if !f(l) {
					return false
				}
			}
			return true
		}
	}
}
