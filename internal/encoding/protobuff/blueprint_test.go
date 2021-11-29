package protobuff

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/blueprint"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/test"
	"testing"
)

func TestMarshalBlueprint(t *testing.T) {
	tests := []*blueprint.Blueprint{
		{
			Name:   blueprintName,
			Stages: []string{bpStage1, bpStage2, bpStage3},
			Links:  []string{bpLink1, bpLink2, bpLink3},
		},
		{
			Name:   "",
			Stages: nil,
			Links:  nil,
		},
	}

	for _, bp := range tests {
		testName := fmt.Sprintf("blueprint=%v", bp)

		t.Run(
			testName, func(t *testing.T) {
				res := MarshalBlueprint(bp)
				test.DeepEqual(t, bp.Name, res.Name, "blueprint name")
				test.DeepEqual(
					t,
					len(bp.Stages),
					len(res.Stages),
					"blueprint stages len")
				for i, s := range bp.Stages {
					test.DeepEqual(t, s, res.Stages[i], "stage %d equal", i)
				}
				test.DeepEqual(
					t,
					len(bp.Links),
					len(res.Links),
					"blueprint links len")
				for i, l := range bp.Links {
					test.DeepEqual(t, l, res.Links[i], "link %d equal", i)
				}
			})
	}
}

func TestUnmarshalBlueprintCorrect(t *testing.T) {
	tests := []*pb.Blueprint{
		{
			Name:   blueprintName,
			Stages: []string{bpStage1, bpStage2, bpStage3},
			Links:  []string{bpLink1, bpLink2, bpLink3},
		},
		{
			Name:   "",
			Stages: nil,
			Links:  nil,
		},
	}

	for _, bp := range tests {
		testName := fmt.Sprintf("blueprint=%v", bp)

		t.Run(
			testName,
			func(t *testing.T) {
				res, err := UnmarshalBlueprint(bp)
				test.DeepEqual(t, nil, err, "unmarshal error")
				test.DeepEqual(t, bp.Name, res.Name, "blueprint name")
				test.DeepEqual(
					t,
					len(bp.Stages),
					len(res.Stages),
					"blueprint stages len")
				for i, s := range bp.Stages {
					test.DeepEqual(t, s, res.Stages[i], "stage %d equal", i)
				}
				test.DeepEqual(
					t,
					len(bp.Links),
					len(res.Links),
					"blueprint links len")
				for i, l := range bp.Links {
					test.DeepEqual(t, l, res.Links[i], "link %d equal", i)
				}
			})
	}
}

func TestUnmarshalBlueprintError(t *testing.T) {
	res, err := UnmarshalBlueprint(nil)
	expectedErr := errdefs.InvalidArgumentWithMsg("'p' is nil")
	test.DeepEqual(t, expectedErr, err, "unmarshal error")
	test.IsNil(t, res, "nil return value")
}
