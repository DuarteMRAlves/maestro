package blueprint

import (
	"github.com/DuarteMRAlves/maestro/internal/assert"
	"testing"
)

const (
	oldBpName = "Old Blueprint Name"
	newBpName = "New Blueprint Name"
)

var (
	stage1 = "Stage 1"
	stage2 = "Stage 2"
	link1  = &Link{SourceField: "Source Field 1"}
	link2  = &Link{TargetField: "Target Field 2"}
)

func TestBlueprint_Clone(t *testing.T) {
	bp := &Blueprint{
		Name:   oldBpName,
		Stages: []string{stage1},
		Links:  []*Link{link1},
	}
	c := bp.Clone()

	assert.DeepEqual(t, oldBpName, c.Name, "cloned old name")
	assert.DeepEqual(t, 1, len(c.Stages), "cloned old Stages length")
	assert.DeepEqual(t, "Stage 1", c.Stages[0], "cloned old stage name")
	assert.DeepEqual(t, 1, len(c.Links), "cloned old Links length")
	assert.DeepEqual(
		t,
		"Source Field 1",
		c.Links[0].SourceField,
		"cloned old link source field")

	c.Name = newBpName
	c.Stages = append(c.Stages, stage2)
	c.Links = append(c.Links, link2)

	assert.DeepEqual(t, oldBpName, bp.Name, "source old name")
	assert.DeepEqual(t, 1, len(bp.Stages), "source old Stages length")
	assert.DeepEqual(t, "Stage 1", bp.Stages[0], "source old stage name")
	assert.DeepEqual(t, 1, len(bp.Links), "source old Links length")
	assert.DeepEqual(
		t,
		"Source Field 1",
		bp.Links[0].SourceField,
		"source old link source field")

	assert.DeepEqual(t, newBpName, c.Name, "cloned new name")
	assert.DeepEqual(t, 2, len(c.Stages), "cloned new Stages length")
	assert.DeepEqual(t, "Stage 1", c.Stages[0], "cloned new stage 1 name")
	assert.DeepEqual(t, "Stage 2", c.Stages[1], "cloned new stage 2 name")
	assert.DeepEqual(t, 2, len(c.Links), "cloned new Links length")
	assert.DeepEqual(
		t,
		"Source Field 1",
		c.Links[0].SourceField,
		"cloned new link 1 source field")
	assert.DeepEqual(
		t,
		"Target Field 2",
		c.Links[1].TargetField,
		"cloned new link 2 target field")
}
