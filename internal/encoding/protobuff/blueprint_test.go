package protobuff

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/blueprint"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"gotest.tools/v3/assert"
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
				assert.Equal(t, bp.Name, res.Name, "blueprint name")
				assert.Equal(
					t,
					len(bp.Stages),
					len(res.Stages),
					"blueprint stages len")
				for i, s := range bp.Stages {
					assert.Equal(t, s, res.Stages[i], "stage %d equal", i)
				}
				assert.Equal(
					t,
					len(bp.Links),
					len(res.Links),
					"blueprint links len")
				for i, l := range bp.Links {
					assert.Equal(t, l, res.Links[i], "link %d equal", i)
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
				assert.Equal(t, nil, err, "unmarshal error")
				assert.Equal(t, bp.Name, res.Name, "blueprint name")
				assert.Equal(
					t,
					len(bp.Stages),
					len(res.Stages),
					"blueprint stages len")
				for i, s := range bp.Stages {
					assert.Equal(t, s, res.Stages[i], "stage %d equal", i)
				}
				assert.Equal(
					t,
					len(bp.Links),
					len(res.Links),
					"blueprint links len")
				for i, l := range bp.Links {
					assert.Equal(t, l, res.Links[i], "link %d equal", i)
				}
			})
	}
}

func TestUnmarshalBlueprintError(t *testing.T) {
	res, err := UnmarshalBlueprint(nil)
	assert.Assert(t, errdefs.IsInvalidArgument(err), "err type")
	assert.ErrorContains(t, err, "'p' is nil")
	assert.Assert(t, res == nil, "nil return value")
}
