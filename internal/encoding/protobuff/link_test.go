package protobuff

import (
	"errors"
	"fmt"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/assert"
	"github.com/DuarteMRAlves/maestro/internal/blueprint"
	"testing"
)

func TestMarshalLink(t *testing.T) {
	tests := []*blueprint.Link{
		{
			SourceStage: linkSourceStage,
			SourceField: linkSourceField,
			TargetStage: linkTargetStage,
			TargetField: linkTargetField,
		},
		{
			SourceStage: "",
			SourceField: "",
			TargetStage: "",
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
			SourceStage: linkSourceStage,
			SourceField: linkSourceField,
			TargetStage: linkTargetStage,
			TargetField: linkTargetField,
		},
		{
			SourceStage: "",
			SourceField: "",
			TargetStage: "",
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
	assert.DeepEqual(t, expected.SourceStage, actual.SourceStage, "source id")
	assert.DeepEqual(
		t,
		expected.SourceField,
		actual.SourceField,
		"source field")
	assert.DeepEqual(t, expected.TargetStage, actual.TargetStage, "target id")
	assert.DeepEqual(
		t,
		expected.TargetField,
		actual.TargetField,
		"target field")
}

func assertPbLink(t *testing.T, expected *pb.Link, actual *blueprint.Link) {
	assert.DeepEqual(t, expected.SourceStage, actual.SourceStage, "source id")
	assert.DeepEqual(
		t,
		expected.SourceField,
		actual.SourceField,
		"source field")
	assert.DeepEqual(t, expected.TargetStage, actual.TargetStage, "target id")
	assert.DeepEqual(
		t,
		expected.TargetField,
		actual.TargetField,
		"target field")
}
