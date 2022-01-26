package flow

import (
	"github.com/DuarteMRAlves/maestro/internal/api"
)

// Flow executes an orchestration.
type Flow struct {
	orchestration *api.Orchestration
}

func New(o *api.Orchestration) *Flow {
	return &Flow{orchestration: o}
}
