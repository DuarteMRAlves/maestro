package protobuff

import (
	"errors"
	"fmt"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/assert"
	"github.com/DuarteMRAlves/maestro/internal/blueprint"
	"github.com/DuarteMRAlves/maestro/internal/identifier"
	"testing"
)

func TestMarshalLink(t *testing.T) {
	sourceId, err := identifier.Rand(5)
	assert.IsNil(t, err, "create source id")
	targetId, err := identifier.Rand(5)
	assert.IsNil(t, err, "create target id")

	tests := []*blueprint.Link{
		{
			SourceId:    sourceId,
			SourceField: linkSourceField,
			TargetId:    targetId,
			TargetField: linkTargetField,
		},
		{
			SourceId:    identifier.Empty(),
			SourceField: "",
			TargetId:    identifier.Empty(),
			TargetField: "",
		},
	}

	for _, l := range tests {
		testName := fmt.Sprintf("link=%v", l)

		t.Run(
			testName, func(t *testing.T) {
				res := MarshalLink(l)
				assertLink(t, l, res)
			})
	}
}

func TestUnmarshalLinkCorrect(t *testing.T) {
	tests := []*pb.Link{
		{
			SourceId:    &pb.Id{Val: linkSourceId},
			SourceField: linkSourceField,
			TargetId:    &pb.Id{Val: linkTargetId},
			TargetField: linkTargetField,
		},
		{
			SourceId:    &pb.Id{Val: ""},
			SourceField: "",
			TargetId:    &pb.Id{Val: ""},
			TargetField: "",
		},
	}
	for _, l := range tests {
		testName := fmt.Sprintf("link=%v", l)

		t.Run(
			testName,
			func(t *testing.T) {
				res, err := UnmarshalLink(l)
				assert.DeepEqual(t, nil, err, "unmarshal error")
				assertPbLink(t, l, res)
			})
	}
}

func TestUnmarshalLinkError(t *testing.T) {
	res, err := UnmarshalLink(nil)
	assert.DeepEqual(t, errors.New("p is nil"), err, "unmarshal error")
	assert.IsNil(t, res, "nil return value")
}

func assertLink(t *testing.T, expected *blueprint.Link, actual *pb.Link) {
	assert.DeepEqual(t, expected.SourceId.Val, actual.SourceId.Val, "source id")
	assert.DeepEqual(
		t,
		expected.SourceField,
		actual.SourceField,
		"source field")
	assert.DeepEqual(t, expected.TargetId.Val, actual.TargetId.Val, "target id")
	assert.DeepEqual(
		t,
		expected.TargetField,
		actual.TargetField,
		"target field")
}

func assertPbLink(t *testing.T, expected *pb.Link, actual *blueprint.Link) {
	assert.DeepEqual(t, expected.SourceId.Val, actual.SourceId.Val, "source id")
	assert.DeepEqual(
		t,
		expected.SourceField,
		actual.SourceField,
		"source field")
	assert.DeepEqual(t, expected.TargetId.Val, actual.TargetId.Val, "target id")
	assert.DeepEqual(
		t,
		expected.TargetField,
		actual.TargetField,
		"target field")
}
