package mapstore

import "github.com/DuarteMRAlves/maestro/internal"

type Stages map[internal.StageName]internal.Stage

func (s Stages) Save(o internal.Stage) error {
	s[o.Name()] = o
	return nil
}

func (s Stages) Load(n internal.StageName) (internal.Stage, error) {
	o, exists := s[n]
	if !exists {
		err := &internal.NotFound{Type: "stage", Ident: n.Unwrap()}
		return internal.Stage{}, err
	}
	return o, nil
}
