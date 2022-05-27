package mapstore

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal"
)

type Links map[internal.LinkName]internal.Link

func (s Links) Save(o internal.Link) error {
	s[o.Name()] = o
	return nil
}

func (s Links) Load(n internal.LinkName) (internal.Link, error) {
	o, exists := s[n]
	if !exists {
		return internal.Link{}, &linkNotFound{name: n.Unwrap()}
	}
	return o, nil
}

type linkNotFound struct{ name string }

func (err *linkNotFound) NotFound() {}

func (err *linkNotFound) Error() string {
	return fmt.Sprintf("link not found: %s", err.name)
}
