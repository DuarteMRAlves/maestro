package flow

import "github.com/DuarteMRAlves/maestro/internal/storage"

// Flow executes an orchestration.
type Flow struct {
	orchestration *storage.Orchestration
}

func New(o *storage.Orchestration) *Flow {
	return &Flow{orchestration: o}
}
