package orchestration

import (
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/validate"
	"sync"
)

// Store manages the storage of created orchestrations.
type Store interface {
	Create(b *Orchestration) error
	Get(query *apitypes.Orchestration) []*Orchestration
}

type store struct {
	orchestrations sync.Map
}

func NewStore() Store {
	return &store{orchestrations: sync.Map{}}
}

func (st *store) Create(config *Orchestration) error {
	if ok, err := validate.ArgNotNil(config, "config"); !ok {
		return err
	}

	o := config.Clone()
	_, prev := st.orchestrations.LoadOrStore(o.Name(), o)
	if prev {
		return errdefs.AlreadyExistsWithMsg(
			"orchestration '%v' already exists",
			o.Name())
	}
	return nil
}

// Get retrieves copies of the stored orchestrations that match the received query.
// The query is a orchestration with the fields that the returned orchestrations should
// have. If a field is empty, then all values for that field are accepted.
func (st *store) Get(query *apitypes.Orchestration) []*Orchestration {
	if query == nil {
		query = &apitypes.Orchestration{}
	}
	filter := buildQueryFilter(query)
	res := make([]*Orchestration, 0)
	st.orchestrations.Range(
		func(key, value interface{}) bool {
			b, ok := value.(*Orchestration)
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

func buildQueryFilter(query *apitypes.Orchestration) func(b *Orchestration) bool {
	filters := make([]func(b *Orchestration) bool, 0)
	if query.Name != "" {
		filters = append(
			filters,
			func(b *Orchestration) bool {
				return b.Name() == query.Name
			})
	}
	if query.Phase != "" {
		filters = append(
			filters,
			func(o *Orchestration) bool {
				return o.phase == query.Phase
			})
	}
	switch len(filters) {
	case 0:
		return func(b *Orchestration) bool { return true }
	case 1:
		return filters[0]
	default:
		return func(b *Orchestration) bool {
			for _, f := range filters {
				if !f(b) {
					return false
				}
			}
			return true
		}
	}
}
