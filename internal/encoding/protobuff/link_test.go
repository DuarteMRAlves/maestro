package protobuff

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/api/pb"
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"gotest.tools/v3/assert"
	"testing"
)

func TestMarshalLink(t *testing.T) {
	const (
		name                           = "name"
		sourceStage apitypes.StageName = "sourceStage"
		sourceField                    = "sourceField"
		targetStage apitypes.StageName = "targetStage"
		targetField                    = "targetField"
	)
	tests := []*apitypes.Link{
		{
			Name:        name,
			SourceStage: sourceStage,
			SourceField: sourceField,
			TargetStage: targetStage,
			TargetField: targetField,
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
				res, err := MarshalLink(l)
				assert.NilError(t, err, "marshal error")
				assertLink(t, l, res)
			})
	}
}

func TestUnmarshalLinkCorrect(t *testing.T) {
	const (
		name        = "name"
		sourceStage = "sourceStage"
		sourceField = "sourceField"
		targetStage = "targetStage"
		targetField = "targetField"
	)
	tests := []*pb.Link{
		{
			Name:        name,
			SourceStage: sourceStage,
			SourceField: sourceField,
			TargetStage: targetStage,
			TargetField: targetField,
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

func TestMarshalLinkNil(t *testing.T) {
	res, err := MarshalLink(nil)
	assert.Assert(t, errdefs.IsInvalidArgument(err), "err type")
	assert.ErrorContains(t, err, "'l' is nil")
	assert.Assert(t, res == nil, "nil return value")
}

func TestUnmarshalLinkNil(t *testing.T) {
	res, err := UnmarshalLink(nil)
	assert.Assert(t, errdefs.IsInvalidArgument(err), "err type")
	assert.ErrorContains(t, err, "'p' is nil")
	assert.Assert(t, res == nil, "nil return value")
}

func assertLink(t *testing.T, expected *apitypes.Link, actual *pb.Link) {
	assert.Equal(t, expected.Name, actual.Name, "name")
	assert.Equal(t, string(expected.SourceStage), actual.SourceStage, "source id")
	assert.Equal(
		t,
		expected.SourceField,
		actual.SourceField,
		"source field")
	assert.Equal(t, string(expected.TargetStage), actual.TargetStage, "target id")
	assert.Equal(
		t,
		expected.TargetField,
		actual.TargetField,
		"target field")
}

func assertPbLink(t *testing.T, expected *pb.Link, actual *apitypes.Link) {
	assert.Equal(t, expected.Name, actual.Name, "name")
	assert.Equal(t, expected.SourceStage, string(actual.SourceStage), "source id")
	assert.Equal(
		t,
		expected.SourceField,
		actual.SourceField,
		"source field")
	assert.Equal(t, expected.TargetStage, string(actual.TargetStage), "target id")
	assert.Equal(
		t,
		expected.TargetField,
		actual.TargetField,
		"target field")
}
