package link

import (
	"github.com/DuarteMRAlves/maestro/internal/test"
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
	test.DeepEqual(t, oldName, c.Name, "cloned old name")
	test.DeepEqual(
		t,
		oldSourceStage,
		c.SourceStage,
		"cloned old source stage")
	test.DeepEqual(
		t,
		oldSourceField,
		c.SourceField,
		"cloned old source field")
	test.DeepEqual(
		t,
		oldTargetStage,
		c.TargetStage,
		"cloned old target stage")
	test.DeepEqual(
		t,
		oldTargetField,
		c.TargetField,
		"cloned old target field")

	c.Name = newName
	c.SourceStage = newSourceStage
	c.SourceField = newSourceField
	c.TargetStage = newTargetStage
	c.TargetField = newTargetField

	test.DeepEqual(t, oldName, s.Name, "source old name")
	test.DeepEqual(
		t,
		oldSourceStage,
		s.SourceStage,
		"source old source stage")
	test.DeepEqual(
		t,
		oldSourceField,
		s.SourceField,
		"source old source field")
	test.DeepEqual(
		t,
		oldTargetStage,
		s.TargetStage,
		"source old target stage")
	test.DeepEqual(
		t,
		oldTargetField,
		s.TargetField,
		"source old target field")

	test.DeepEqual(t, newName, c.Name, "cloned new name")
	test.DeepEqual(
		t,
		newSourceStage,
		c.SourceStage,
		"cloned new source stage")
	test.DeepEqual(
		t,
		newSourceField,
		c.SourceField,
		"cloned new source field")
	test.DeepEqual(
		t,
		newTargetStage,
		c.TargetStage,
		"cloned new target stage")
	test.DeepEqual(
		t,
		newTargetField,
		c.TargetField,
		"cloned new target field")
}
