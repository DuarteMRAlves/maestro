package blueprint

import (
	"gotest.tools/v3/assert"
	"testing"
)

const (
	oldBpName = "Old Blueprint Name"
	newBpName = "New Blueprint Name"
	link1     = "Link 1"
	link2     = "Link 2"
)

func TestBlueprint_Clone(t *testing.T) {
	bp := &Blueprint{
		Name:  oldBpName,
		Links: []string{link1},
	}
	c := bp.Clone()

	assert.Equal(t, oldBpName, c.Name, "cloned old name")
	assert.Equal(t, 1, len(c.Links), "cloned old Links length")
	assert.Equal(t, link1, c.Links[0], "cloned old link")

	c.Name = newBpName
	c.Links = append(c.Links, link2)

	assert.Equal(t, oldBpName, bp.Name, "source old name")
	assert.Equal(t, 1, len(bp.Links), "source old Links length")
	assert.Equal(t, link1, bp.Links[0], "source old link name")

	assert.Equal(t, newBpName, c.Name, "cloned new name")
	assert.Equal(t, 2, len(c.Links), "cloned new Links length")
	assert.Equal(t, link1, c.Links[0], "cloned new link 1")
	assert.Equal(t, link2, c.Links[1], "cloned new link 2")
}
