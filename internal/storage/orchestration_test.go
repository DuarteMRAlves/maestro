package storage

import (
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"gotest.tools/v3/assert"
	"testing"
)

func TestOrchestration_Clone(t *testing.T) {
	const (
		oldName  apitypes.OrchestrationName = "Old Orchestration Name"
		oldPhase                            = apitypes.OrchestrationRunning
		newName  apitypes.OrchestrationName = "New Orchestration Name"
		newPhase                            = apitypes.OrchestrationFailed
	)
	o := &Orchestration{
		name:  oldName,
		phase: oldPhase,
	}
	c := o.Clone()

	assert.Equal(t, oldName, c.name, "cloned old name")
	assert.Equal(t, oldPhase, c.phase, "cloned old phase")

	c.name = newName
	c.phase = newPhase

	assert.Equal(t, oldName, o.name, "source old name")
	assert.Equal(t, oldPhase, o.phase, "source old phase")

	assert.Equal(t, newName, c.name, "cloned new name")
	assert.Equal(t, newPhase, c.phase, "cloned new phase")
}
