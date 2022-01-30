package execution

import (
	"github.com/DuarteMRAlves/maestro/internal/api"
)

// Flow executes an orchestration.
type Flow struct {
	orchestration *api.Orchestration
}

func NewFlow(o *api.Orchestration) *Flow {
	return &Flow{orchestration: o}
}
