package asset

import (
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/validate"
	"sync"
)

type Store interface {
	Create(description *Asset) error
	Contains(name string) bool
	Get(query *apitypes.Asset) []*Asset
}

type store struct {
	assets sync.Map
}

func NewStore() Store {
	return &store{assets: sync.Map{}}
}

func (st *store) Create(config *Asset) error {
	if ok, err := validate.ArgNotNil(config, "config"); !ok {
		return errdefs.InvalidArgumentWithError(err)
	}

	asset := config.Clone()
	_, prev := st.assets.LoadOrStore(asset.Name, asset)
	if prev {
		return errdefs.AlreadyExistsWithMsg(
			"asset '%v' already exists",
			asset.Name)
	}
	return nil
}

func (st *store) Contains(name string) bool {
	_, ok := st.assets.Load(name)
	return ok
}

func (st *store) Get(query *apitypes.Asset) []*Asset {
	if query == nil {
		query = &apitypes.Asset{}
	}
	filter := buildQueryFilter(query)
	res := make([]*Asset, 0)
	st.assets.Range(
		func(key, value interface{}) bool {
			a, ok := value.(*Asset)
			if !ok {
				return false
			}
			if filter(a) {
				res = append(res, a.Clone())
			}
			return true
		})
	return res
}

func buildQueryFilter(query *apitypes.Asset) func(a *Asset) bool {
	filters := make([]func(a *Asset) bool, 0)
	if query.Name != "" {
		filters = append(
			filters,
			func(a *Asset) bool {
				return a.Name == query.Name
			})
	}
	if query.Image != "" {
		filters = append(
			filters,
			func(a *Asset) bool {
				return a.Image == query.Image
			})
	}
	if len(filters) > 0 {
		return func(a *Asset) bool {
			for _, f := range filters {
				if !f(a) {
					return false
				}
			}
			return true
		}
	}
	return func(a *Asset) bool {
		return true
	}
}
