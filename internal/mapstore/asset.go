package mapstore

import "github.com/DuarteMRAlves/maestro/internal"

type Assets map[internal.AssetName]internal.Asset

func (s Assets) Save(o internal.Asset) error {
	s[o.Name()] = o
	return nil
}

func (s Assets) Load(n internal.AssetName) (internal.Asset, error) {
	o, exists := s[n]
	if !exists {
		err := &internal.NotFound{Type: "asset", Ident: n.Unwrap()}
		return internal.Asset{}, err
	}
	return o, nil
}
