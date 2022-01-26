package flow

import (
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
)

// Flow executes an orchestration.
type Flow struct {
	orchestration *apitypes.Orchestration
}

func New(o *apitypes.Orchestration) *Flow {
	return &Flow{orchestration: o}
}
