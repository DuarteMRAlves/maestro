package mapstore

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal"
)

type Assets map[internal.AssetName]internal.Asset

func (s Assets) Save(o internal.Asset) error {
	s[o.Name()] = o
	return nil
}

func (s Assets) Load(n internal.AssetName) (internal.Asset, error) {
	o, exists := s[n]
	if !exists {
		return internal.Asset{}, &assetNotFound{name: n.Unwrap()}
	}
	return o, nil
}

type assetNotFound struct{ name string }

func (err *assetNotFound) NotFound() {}

func (err *assetNotFound) Error() string {
	return fmt.Sprintf("asset not found: %s", err.name)
}
