package protobuff

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/link"
	"github.com/DuarteMRAlves/maestro/internal/test"
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
				test.DeepEqual(t, nil, err, "unmarshal error")
				assertPbLink(t, l, res)
			})
	}
}

func TestUnmarshalLinkError(t *testing.T) {
	res, err := UnmarshalLink(nil)
	expectedErr := errdefs.InvalidArgumentWithMsg("'p' is nil")
	test.DeepEqual(t, expectedErr, err, "unmarshal error")
	test.IsNil(t, res, "nil return value")
}

func assertLink(t *testing.T, expected *link.Link, actual *pb.Link) {
	test.DeepEqual(t, expected.Name, actual.Name, "name")
	test.DeepEqual(t, expected.SourceStage, actual.SourceStage, "source id")
	test.DeepEqual(
		t,
		expected.SourceField,
		actual.SourceField,
		"source field")
	test.DeepEqual(t, expected.TargetStage, actual.TargetStage, "target id")
	test.DeepEqual(
		t,
		expected.TargetField,
		actual.TargetField,
		"target field")
}

func assertPbLink(t *testing.T, expected *pb.Link, actual *link.Link) {
	test.DeepEqual(t, expected.Name, actual.Name, "name")
	test.DeepEqual(t, expected.SourceStage, actual.SourceStage, "source id")
	test.DeepEqual(
		t,
		expected.SourceField,
		actual.SourceField,
		"source field")
	test.DeepEqual(t, expected.TargetStage, actual.TargetStage, "target id")
	test.DeepEqual(
		t,
		expected.TargetField,
		actual.TargetField,
		"target field")
}
