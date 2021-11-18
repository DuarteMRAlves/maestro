package blueprint

import (
	"github.com/DuarteMRAlves/maestro/internal/assert"
	"testing"
)

const (
	oldSourceStage = "OldSourceStage"
	oldSourceField = "OldSourceField"
	oldTargetStage = "OldTargetStage"
	oldTargetField = "OldTargetField"

	newSourceStage = "NewSourceStage"
	newSourceField = "NewSourceField"
	newTargetStage = "NewTargetStage"
	newTargetField = "NewTargetField"
)

func TestLink_Clone(t *testing.T) {
	s := &Link{
		SourceStage: oldSourceStage,
		SourceField: oldSourceField,
		TargetStage: oldTargetStage,
		TargetField: oldTargetField,
	}
	c := s.Clone()
	assert.DeepEqual(
		t,
		oldSourceStage,
		c.SourceStage,
		"cloned old source stage")
	assert.DeepEqual(
		t,
		oldSourceField,
		c.SourceField,
		"cloned old source field")
	assert.DeepEqual(
		t,
		oldTargetStage,
		c.TargetStage,
		"cloned old target stage")
	assert.DeepEqual(
		t,
		oldTargetField,
		c.TargetField,
		"cloned old target field")

	c.SourceStage = newSourceStage
	c.SourceField = newSourceField
	c.TargetStage = newTargetStage
	c.TargetField = newTargetField

	assert.DeepEqual(
		t,
		oldSourceStage,
		s.SourceStage,
		"source old source stage")
	assert.DeepEqual(
		t,
		oldSourceField,
		s.SourceField,
		"source old source field")
	assert.DeepEqual(
		t,
		oldTargetStage,
		s.TargetStage,
		"source old target stage")
	assert.DeepEqual(
		t,
		oldTargetField,
		s.TargetField,
		"source old target field")

	assert.DeepEqual(
		t,
		newSourceStage,
		c.SourceStage,
		"cloned new source stage")
	assert.DeepEqual(
		t,
		newSourceField,
		c.SourceField,
		"cloned new source field")
	assert.DeepEqual(
		t,
		newTargetStage,
		c.TargetStage,
		"cloned new target stage")
	assert.DeepEqual(
		t,
		newTargetField,
		c.TargetField,
		"cloned new target field")
}
