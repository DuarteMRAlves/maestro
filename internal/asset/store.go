package asset

import (
	"errors"
	"sync"
)

type Store interface {
	Create(description *Asset) error
	Get(query *Asset) []*Asset
}

type store struct {
	assets sync.Map
}

func NewStore() Store {
	return &store{assets: sync.Map{}}
}

func (st *store) Create(config *Asset) error {
	if config == nil {
		return errors.New("nil config")
	}

	asset := config.Clone()
	_, prev := st.assets.LoadOrStore(asset.Name, asset)
	if prev {
		return AlreadyExists{Name: asset.Name}
	}
	return nil
}

func (st *store) Get(query *Asset) []*Asset {
	if query == nil {
		query = &Asset{}
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

func buildQueryFilter(query *Asset) func(a *Asset) bool {
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
