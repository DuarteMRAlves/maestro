package blueprint

import (
	testing2 "github.com/DuarteMRAlves/maestro/internal/testing"
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

	testing2.DeepEqual(t, oldBpName, c.Name, "cloned old name")
	testing2.DeepEqual(t, 1, len(c.Stages), "cloned old Stages length")
	testing2.DeepEqual(t, stage1, c.Stages[0], "cloned old stage name")
	testing2.DeepEqual(t, 1, len(c.Links), "cloned old Links length")
	testing2.DeepEqual(t, link1, c.Links[0], "cloned old link")

	c.Name = newBpName
	c.Stages = append(c.Stages, stage2)
	c.Links = append(c.Links, link2)

	testing2.DeepEqual(t, oldBpName, bp.Name, "source old name")
	testing2.DeepEqual(t, 1, len(bp.Stages), "source old Stages length")
	testing2.DeepEqual(t, stage1, bp.Stages[0], "source old stage name")
	testing2.DeepEqual(t, 1, len(bp.Links), "source old Links length")
	testing2.DeepEqual(t, link1, bp.Links[0], "source old link name")

	testing2.DeepEqual(t, newBpName, c.Name, "cloned new name")
	testing2.DeepEqual(t, 2, len(c.Stages), "cloned new Stages length")
	testing2.DeepEqual(t, stage1, c.Stages[0], "cloned new stage 1 name")
	testing2.DeepEqual(t, stage2, c.Stages[1], "cloned new stage 2 name")
	testing2.DeepEqual(t, 2, len(c.Links), "cloned new Links length")
	testing2.DeepEqual(t, link1, c.Links[0], "cloned new link 1")
	testing2.DeepEqual(t, link2, c.Links[1], "cloned new link 2")
}
