package link

import (
	"gotest.tools/v3/assert"
	"testing"
)

func TestLink_Clone(t *testing.T) {
	const (
		oldName        = "OldName"
		oldSourceStage = "OldSourceStage"
		oldSourceField = "OldSourceField"
		oldTargetStage = "OldTargetStage"
		oldTargetField = "OldTargetField"

		newName        = "NewName"
		newSourceStage = "NewSourceStage"
		newSourceField = "NewSourceField"
		newTargetStage = "NewTargetStage"
		newTargetField = "NewTargetField"
	)
	s := &Link{
		name:        oldName,
		sourceStage: oldSourceStage,
		sourceField: oldSourceField,
		targetStage: oldTargetStage,
		targetField: oldTargetField,
	}
	c := s.Clone()
	assert.Equal(t, oldName, c.name, "cloned old name")
	assert.Equal(t, oldSourceStage, c.sourceStage, "cloned old source stage")
	assert.Equal(t, oldSourceField, c.sourceField, "cloned old source field")
	assert.Equal(t, oldTargetStage, c.targetStage, "cloned old target stage")
	assert.Equal(t, oldTargetField, c.targetField, "cloned old target field")

	c.name = newName
	c.sourceStage = newSourceStage
	c.sourceField = newSourceField
	c.targetStage = newTargetStage
	c.targetField = newTargetField

	assert.Equal(t, oldName, s.name, "source old name")
	assert.Equal(t, oldSourceStage, s.sourceStage, "source old source stage")
	assert.Equal(t, oldSourceField, s.sourceField, "source old source field")
	assert.Equal(t, oldTargetStage, s.targetStage, "source old target stage")
	assert.Equal(t, oldTargetField, s.targetField, "source old target field")

	assert.Equal(t, newName, c.name, "cloned new name")
	assert.Equal(t, newSourceStage, c.sourceStage, "cloned new source stage")
	assert.Equal(t, newSourceField, c.sourceField, "cloned new source field")
	assert.Equal(t, newTargetStage, c.targetStage, "cloned new target stage")
	assert.Equal(t, newTargetField, c.targetField, "cloned new target field")
}
