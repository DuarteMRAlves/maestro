package execution

import (
	"github.com/DuarteMRAlves/maestro/internal/api"
)

// Execution executes an orchestration.
type Execution struct {
	orchestration *api.Orchestration

	workers map[api.StageName]Worker
}
