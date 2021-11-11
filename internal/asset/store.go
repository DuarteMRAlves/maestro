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
	assets map[identifier.Id]*Asset
	lock   sync.RWMutex
}

func NewStore() Store {
	return &store{assets: map[identifier.Id]*Asset{}, lock: sync.RWMutex{}}
}

func (st *store) Create(config *Asset) (identifier.Id, error) {
	if config == nil {
		return identifier.Empty(), errors.New("nil config")
	}
	if !config.Id.IsEmpty() {
		return identifier.Empty(), errors.New("asset identifier should not be defined")
	}
	asset := config.Clone()

	st.lock.Lock()
	id := st.generateNewId()
	st.assets[id] = asset
	st.lock.Unlock()

	asset.Id = id

	return id.Clone(), nil
}

func (st *store) Get(id identifier.Id) (*Asset, error) {
	if !id.IsValid() {
		return nil, errors.New("invalid identifier")
	}
	st.lock.RLock()
	asset, ok := st.assets[id]
	st.lock.RUnlock()
	if !ok {
		return nil, fmt.Errorf("no such identifier: '%s'", id)
	}
	return asset.Clone(), nil
}

func (st *store) List() ([]*Asset, error) {
	st.lock.RLock()
	defer st.lock.RUnlock()
	res := make([]*Asset, 0, len(st.assets))
	for _, a := range st.assets {
		res = append(res, a.Clone())
	}
	return res, nil
}

func (st *store) generateNewId() identifier.Id {
	newId, err := identifier.Rand(IdSize)
	if err != nil {
		panic(err)
	}
	_, idExists := st.assets[newId]
	for idExists {
		if newId, err = identifier.Rand(IdSize); err != nil {
			panic(err)
		}
		_, idExists = st.assets[newId]
	}
	return newId
}
