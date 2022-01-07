package orchestration

import (
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"gotest.tools/v3/assert"
	"testing"
)

func TestOrchestration_Clone(t *testing.T) {
	const (
		oldBpName apitypes.OrchestrationName = "Old Orchestration Name"
		newBpName apitypes.OrchestrationName = "New Orchestration Name"
		link1                                = "Link 1"
		link2                                = "Link 2"
	)
	o := &Orchestration{
		name:  oldBpName,
		links: []string{link1},
	}
	c := o.Clone()

	assert.Equal(t, oldBpName, c.name, "cloned old name")
	assert.Equal(t, 1, len(c.links), "cloned old Links length")
	assert.Equal(t, link1, c.links[0], "cloned old link")

	c.name = newBpName
	c.links = append(c.links, link2)

	assert.Equal(t, oldBpName, o.name, "source old name")
	assert.Equal(t, 1, len(o.links), "source old Links length")
	assert.Equal(t, link1, o.links[0], "source old link name")

	assert.Equal(t, newBpName, c.name, "cloned new name")
	assert.Equal(t, 2, len(c.links), "cloned new Links length")
	assert.Equal(t, link1, c.links[0], "cloned new link 1")
	assert.Equal(t, link2, c.links[1], "cloned new link 2")
}
