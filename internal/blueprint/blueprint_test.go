package blueprint

import (
	"github.com/DuarteMRAlves/maestro/internal/test"
	"testing"
)

const (
	oldBpName = "Old Blueprint Name"
	newBpName = "New Blueprint Name"
	stage1    = "Stage 1"
	stage2    = "Stage 2"
	link1     = "Link 1"
	link2     = "Link 2"
)

func TestBlueprint_Clone(t *testing.T) {
	bp := &Blueprint{
		Name:   oldBpName,
		Stages: []string{stage1},
		Links:  []string{link1},
	}
	c := bp.Clone()

	test.DeepEqual(t, oldBpName, c.Name, "cloned old name")
	test.DeepEqual(t, 1, len(c.Stages), "cloned old Stages length")
	test.DeepEqual(t, stage1, c.Stages[0], "cloned old stage name")
	test.DeepEqual(t, 1, len(c.Links), "cloned old Links length")
	test.DeepEqual(t, link1, c.Links[0], "cloned old link")

	c.Name = newBpName
	c.Stages = append(c.Stages, stage2)
	c.Links = append(c.Links, link2)

	test.DeepEqual(t, oldBpName, bp.Name, "source old name")
	test.DeepEqual(t, 1, len(bp.Stages), "source old Stages length")
	test.DeepEqual(t, stage1, bp.Stages[0], "source old stage name")
	test.DeepEqual(t, 1, len(bp.Links), "source old Links length")
	test.DeepEqual(t, link1, bp.Links[0], "source old link name")

	test.DeepEqual(t, newBpName, c.Name, "cloned new name")
	test.DeepEqual(t, 2, len(c.Stages), "cloned new Stages length")
	test.DeepEqual(t, stage1, c.Stages[0], "cloned new stage 1 name")
	test.DeepEqual(t, stage2, c.Stages[1], "cloned new stage 2 name")
	test.DeepEqual(t, 2, len(c.Links), "cloned new Links length")
	test.DeepEqual(t, link1, c.Links[0], "cloned new link 1")
	test.DeepEqual(t, link2, c.Links[1], "cloned new link 2")
}
