package link

import (
	"gotest.tools/v3/assert"
	"testing"
)

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

func TestLink_Clone(t *testing.T) {
	s := &Link{
		Name:        oldName,
		SourceStage: oldSourceStage,
		SourceField: oldSourceField,
		TargetStage: oldTargetStage,
		TargetField: oldTargetField,
	}
	c := s.Clone()
	assert.Equal(t, oldName, c.Name, "cloned old name")
	assert.Equal(t, oldSourceStage, c.SourceStage, "cloned old source stage")
	assert.Equal(t, oldSourceField, c.SourceField, "cloned old source field")
	assert.Equal(t, oldTargetStage, c.TargetStage, "cloned old target stage")
	assert.Equal(t, oldTargetField, c.TargetField, "cloned old target field")

	c.Name = newName
	c.SourceStage = newSourceStage
	c.SourceField = newSourceField
	c.TargetStage = newTargetStage
	c.TargetField = newTargetField

	assert.Equal(t, oldName, s.Name, "source old name")
	assert.Equal(t, oldSourceStage, s.SourceStage, "source old source stage")
	assert.Equal(t, oldSourceField, s.SourceField, "source old source field")
	assert.Equal(t, oldTargetStage, s.TargetStage, "source old target stage")
	assert.Equal(t, oldTargetField, s.TargetField, "source old target field")

	assert.Equal(t, newName, c.Name, "cloned new name")
	assert.Equal(t, newSourceStage, c.SourceStage, "cloned new source stage")
	assert.Equal(t, newSourceField, c.SourceField, "cloned new source field")
	assert.Equal(t, newTargetStage, c.TargetStage, "cloned new target stage")
	assert.Equal(t, newTargetField, c.TargetField, "cloned new target field")
}
