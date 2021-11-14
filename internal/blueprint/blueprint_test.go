package blueprint

import (
	"github.com/DuarteMRAlves/maestro/internal/assert"
	"github.com/DuarteMRAlves/maestro/internal/identifier"
	"testing"
)

const (
	oldBpId   = "AGR423FE53"
	oldBpName = "Old Blueprint Name"
	newBpId   = "1RGN92ND98"
	newBpName = "New Blueprint Name"
)

var (
	stage1 = &Stage{Name: "Stage 1"}
	stage2 = &Stage{Name: "Stage 2"}
	link1  = &Link{SourceField: "Source Field 1"}
	link2  = &Link{TargetField: "Target Field 2"}
)

func TestBlueprint_Clone(t *testing.T) {
	bp := &Blueprint{
		Id:     identifier.Id{Val: oldBpId},
		Name:   oldBpName,
		stages: []*Stage{stage1},
		links:  []*Link{link1},
	}
	c := bp.Clone()

	assert.DeepEqual(t, oldBpId, c.Id.Val, "cloned old id")
	assert.DeepEqual(t, oldBpName, c.Name, "cloned old name")
	assert.DeepEqual(t, 1, len(c.stages), "cloned old stages length")
	assert.DeepEqual(t, "Stage 1", c.stages[0].Name, "cloned old stage name")
	assert.DeepEqual(t, 1, len(c.links), "cloned old links length")
	assert.DeepEqual(
		t,
		"Source Field 1",
		c.links[0].SourceField,
		"cloned old link source field")

	c.Id.Val = newBpId
	c.Name = newBpName
	c.stages = append(c.stages, stage2)
	c.links = append(c.links, link2)

	assert.DeepEqual(t, oldBpId, bp.Id.Val, "source old id")
	assert.DeepEqual(t, oldBpName, bp.Name, "source old name")
	assert.DeepEqual(t, 1, len(bp.stages), "source old stages length")
	assert.DeepEqual(t, "Stage 1", bp.stages[0].Name, "source old stage name")
	assert.DeepEqual(t, 1, len(bp.links), "source old links length")
	assert.DeepEqual(
		t,
		"Source Field 1",
		bp.links[0].SourceField,
		"source old link source field")

	assert.DeepEqual(t, newBpId, c.Id.Val, "cloned new id")
	assert.DeepEqual(t, newBpName, c.Name, "cloned new name")
	assert.DeepEqual(t, 2, len(c.stages), "cloned new stages length")
	assert.DeepEqual(t, "Stage 1", c.stages[0].Name, "cloned new stage 1 name")
	assert.DeepEqual(t, "Stage 2", c.stages[1].Name, "cloned new stage 2 name")
	assert.DeepEqual(t, 2, len(c.links), "cloned new links length")
	assert.DeepEqual(
		t,
		"Source Field 1",
		c.links[0].SourceField,
		"cloned new link 1 source field")
	assert.DeepEqual(
		t,
		"Target Field 2",
		c.links[1].TargetField,
		"cloned new link 2 target field")
}
