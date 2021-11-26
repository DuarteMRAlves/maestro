package link

import (
	testing2 "github.com/DuarteMRAlves/maestro/internal/testing"
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
	testing2.DeepEqual(t, oldName, c.Name, "cloned old name")
	testing2.DeepEqual(
		t,
		oldSourceStage,
		c.SourceStage,
		"cloned old source stage")
	testing2.DeepEqual(
		t,
		oldSourceField,
		c.SourceField,
		"cloned old source field")
	testing2.DeepEqual(
		t,
		oldTargetStage,
		c.TargetStage,
		"cloned old target stage")
	testing2.DeepEqual(
		t,
		oldTargetField,
		c.TargetField,
		"cloned old target field")

	c.Name = newName
	c.SourceStage = newSourceStage
	c.SourceField = newSourceField
	c.TargetStage = newTargetStage
	c.TargetField = newTargetField

	testing2.DeepEqual(t, oldName, s.Name, "source old name")
	testing2.DeepEqual(
		t,
		oldSourceStage,
		s.SourceStage,
		"source old source stage")
	testing2.DeepEqual(
		t,
		oldSourceField,
		s.SourceField,
		"source old source field")
	testing2.DeepEqual(
		t,
		oldTargetStage,
		s.TargetStage,
		"source old target stage")
	testing2.DeepEqual(
		t,
		oldTargetField,
		s.TargetField,
		"source old target field")

	testing2.DeepEqual(t, newName, c.Name, "cloned new name")
	testing2.DeepEqual(
		t,
		newSourceStage,
		c.SourceStage,
		"cloned new source stage")
	testing2.DeepEqual(
		t,
		newSourceField,
		c.SourceField,
		"cloned new source field")
	testing2.DeepEqual(
		t,
		newTargetStage,
		c.TargetStage,
		"cloned new target stage")
	testing2.DeepEqual(
		t,
		newTargetField,
		c.TargetField,
		"cloned new target field")
}
