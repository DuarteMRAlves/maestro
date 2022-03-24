package mapstore

import "github.com/DuarteMRAlves/maestro/internal"

type Orchestrations map[internal.OrchestrationName]internal.Orchestration

func (s Orchestrations) Save(o internal.Orchestration) error {
	s[o.Name()] = o
	return nil
}

func (s Orchestrations) Load(
	n internal.OrchestrationName,
) (internal.Orchestration, error) {
	o, exists := s[n]
	if !exists {
		err := &internal.NotFound{Type: "orchestration", Ident: n.Unwrap()}
		return internal.Orchestration{}, err
	}
	return o, nil
}
