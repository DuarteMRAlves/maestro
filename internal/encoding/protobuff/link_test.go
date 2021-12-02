package protobuff

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/link"
	"gotest.tools/v3/assert"
	"testing"
)

func TestMarshalLink(t *testing.T) {
	tests := []*link.Link{
		{
			Name:        linkName,
			SourceStage: linkSourceStage,
			SourceField: linkSourceField,
			TargetStage: linkTargetStage,
			TargetField: linkTargetField,
		},
		{
			Name:        "",
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
			Name:        linkName,
			SourceStage: linkSourceStage,
			SourceField: linkSourceField,
			TargetStage: linkTargetStage,
			TargetField: linkTargetField,
		},
		{
			Name:        "",
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
				assert.NilError(t, err, "unmarshal error")
				assertPbLink(t, l, res)
			})
	}
}

func TestUnmarshalLinkError(t *testing.T) {
	res, err := UnmarshalLink(nil)
	assert.Assert(t, errdefs.IsInvalidArgument(err), "err type")
	assert.ErrorContains(t, err, "'p' is nil")
	assert.Assert(t, res == nil, "nil return value")
}

func assertLink(t *testing.T, expected *link.Link, actual *pb.Link) {
	assert.Equal(t, expected.Name, actual.Name, "name")
	assert.Equal(t, expected.SourceStage, actual.SourceStage, "source id")
	assert.Equal(
		t,
		expected.SourceField,
		actual.SourceField,
		"source field")
	assert.Equal(t, expected.TargetStage, actual.TargetStage, "target id")
	assert.Equal(
		t,
		expected.TargetField,
		actual.TargetField,
		"target field")
}

func assertPbLink(t *testing.T, expected *pb.Link, actual *link.Link) {
	assert.Equal(t, expected.Name, actual.Name, "name")
	assert.Equal(t, expected.SourceStage, actual.SourceStage, "source id")
	assert.Equal(
		t,
		expected.SourceField,
		actual.SourceField,
		"source field")
	assert.Equal(t, expected.TargetStage, actual.TargetStage, "target id")
	assert.Equal(
		t,
		expected.TargetField,
		actual.TargetField,
		"target field")
}
