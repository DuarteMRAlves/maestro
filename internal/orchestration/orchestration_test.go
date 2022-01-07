package orchestration

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
		link1    apitypes.LinkName          = "Link 1"
		link2    apitypes.LinkName          = "Link 2"
	)
	o := &Orchestration{
		name:  oldName,
		phase: oldPhase,
		links: []apitypes.LinkName{link1},
	}
	c := o.Clone()

	assert.Equal(t, oldName, c.name, "cloned old name")
	assert.Equal(t, oldPhase, c.phase, "cloned old phase")
	assert.Equal(t, 1, len(c.links), "cloned old Links length")
	assert.Equal(t, link1, c.links[0], "cloned old link")

	c.name = newName
	c.phase = newPhase
	c.links = append(c.links, link2)

	assert.Equal(t, oldName, o.name, "source old name")
	assert.Equal(t, oldPhase, o.phase, "source old phase")
	assert.Equal(t, 1, len(o.links), "source old Links length")
	assert.Equal(t, link1, o.links[0], "source old link name")

	assert.Equal(t, newName, c.name, "cloned new name")
	assert.Equal(t, newPhase, c.phase, "cloned new phase")
	assert.Equal(t, 2, len(c.links), "cloned new Links length")
	assert.Equal(t, link1, c.links[0], "cloned new link 1")
	assert.Equal(t, link2, c.links[1], "cloned new link 2")
}
