package blueprint

import (
	"github.com/DuarteMRAlves/maestro/internal/assert"
	"github.com/DuarteMRAlves/maestro/internal/identifier"
	"testing"
)

const (
	oldSourceId    = "GIB343AFE1"
	oldSourceField = "OldSourceField"
	oldTargetId    = "Y5UYV878BU"
	oldTargetField = "OldTargetField"

	newSourceId    = "GIB564VV97"
	newSourceField = "NewSourceField"
	newTargetId    = "GOEF67V5CD"
	newTargetField = "NewTargetField"
)

func TestLink_Clone(t *testing.T) {
	s := &Link{
		SourceId:    identifier.Id{Val: oldSourceId},
		SourceField: oldSourceField,
		TargetId:    identifier.Id{Val: oldTargetId},
		TargetField: oldTargetField,
	}
	c := s.Clone()
	assert.DeepEqual(t, oldSourceId, c.SourceId.Val, "cloned old source id")
	assert.DeepEqual(
		t,
		oldSourceField,
		c.SourceField,
		"cloned old source field")
	assert.DeepEqual(t, oldTargetId, c.TargetId.Val, "cloned old target id")
	assert.DeepEqual(
		t,
		oldTargetField,
		c.TargetField,
		"cloned old target field")

	c.SourceId.Val = newSourceId
	c.SourceField = newSourceField
	c.TargetId.Val = newTargetId
	c.TargetField = newTargetField

	assert.DeepEqual(t, oldSourceId, s.SourceId.Val, "source old source id")
	assert.DeepEqual(
		t,
		oldSourceField,
		s.SourceField,
		"source old source field")
	assert.DeepEqual(t, oldTargetId, s.TargetId.Val, "source old target id")
	assert.DeepEqual(
		t,
		oldTargetField,
		s.TargetField,
		"source old target field")

	assert.DeepEqual(t, newSourceId, c.SourceId.Val, "cloned new source id")
	assert.DeepEqual(
		t,
		newSourceField,
		c.SourceField,
		"cloned new source field")
	assert.DeepEqual(t, newTargetId, c.TargetId.Val, "cloned new target id")
	assert.DeepEqual(
		t,
		newTargetField,
		c.TargetField,
		"cloned new target field")
}
