package flow

import "github.com/DuarteMRAlves/maestro/internal/orchestration"

// Flow executes an orchestration.
type Flow struct {
	orchestration *orchestration.Orchestration
}

func New(o *orchestration.Orchestration) *Flow {
	return &Flow{orchestration: o}
}
