package asset

import (
	"errors"
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/identifier"
	"sync"
)

type Store interface {
	Create(description *Asset) (identifier.Id, error)
	Get(id identifier.Id) (*Asset, error)
	List() ([]*Asset, error)
}

type store struct {
	assets sync.Map
	gen    identifier.Generator
}

func NewStore() Store {
	return &store{
		assets: sync.Map{},
		gen:    identifier.GenForSize(IdSize),
	}
}

func (st *store) Create(config *Asset) (identifier.Id, error) {
	if config == nil {
		return identifier.Empty(), errors.New("nil config")
	}
	if !config.Id.IsEmpty() {
		return identifier.Empty(), errors.New("asset identifier should not be defined")
	}
	asset := config.Clone()

	id, err := st.gen()
	if err != nil {
		return identifier.Empty(), fmt.Errorf("store create asset: %v", err)
	}
	st.assets.Store(id, asset)

	asset.Id = id

	return id.Clone(), nil
}

func (st *store) Get(id identifier.Id) (*Asset, error) {
	if !id.IsValid() {
		return nil, errors.New("invalid identifier")
	}

	stored, ok := st.assets.Load(id)
	if !ok {
		return nil, fmt.Errorf("no such identifier: '%s'", id)
	}
	asset, ok := stored.(*Asset)
	if !ok {
		return nil, fmt.Errorf("store get asset %v: type assertion failed", id)
	}
	return asset.Clone(), nil
}

func (st *store) List() ([]*Asset, error) {
	res := make([]*Asset, 0)
	st.assets.Range(
		func(key, value interface{}) bool {
			a, ok := value.(*Asset)
			if !ok {
				return false
			}
			res = append(res, a.Clone())
			return true
		})
	return res, nil
}
