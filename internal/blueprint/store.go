package blueprint

import (
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/validate"
	"sync"
)

// Store manages the storage of created blueprints.
type Store interface {
	Create(b *Blueprint) error
	Get(query *Blueprint) []*Blueprint
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

// Get retrieves copies of the stored blueprints that match the received query.
// The query is a blueprint with the fields that the returned blueprints should
// have. If a field is empty, then all values for that field are accepted.
func (st *store) Get(query *Blueprint) []*Blueprint {
	if query == nil {
		query = &Blueprint{}
	}
	filter := buildQueryFilter(query)
	res := make([]*Blueprint, 0)
	st.blueprints.Range(
		func(key, value interface{}) bool {
			b, ok := value.(*Blueprint)
			if !ok {
				return false
			}
			if filter(b) {
				res = append(res, b.Clone())
			}
			return true
		})
	return res
}

func buildQueryFilter(query *Blueprint) func(b *Blueprint) bool {
	filters := make([]func(b *Blueprint) bool, 0)
	if query.Name != "" {
		filters = append(
			filters,
			func(b *Blueprint) bool {
				return b.Name == query.Name
			})
	}
	switch len(filters) {
	case 0:
		return func(b *Blueprint) bool { return true }
	case 1:
		return filters[0]
	default:
		return func(b *Blueprint) bool {
			for _, f := range filters {
				if !f(b) {
					return false
				}
			}
			return true
		}
	}
}
