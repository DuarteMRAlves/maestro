package mapstore

import "github.com/DuarteMRAlves/maestro/internal"

type Links map[internal.LinkName]internal.Link

func (s Links) Save(o internal.Link) error {
	s[o.Name()] = o
	return nil
}

func (s Links) Load(n internal.LinkName) (internal.Link, error) {
	o, exists := s[n]
	if !exists {
		err := &internal.NotFound{Type: "link", Ident: n.Unwrap()}
		return internal.Link{}, err
	}
	return o, nil
}
