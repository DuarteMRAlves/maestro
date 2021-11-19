package protobuff

import (
	"errors"
	"fmt"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/assert"
	"github.com/DuarteMRAlves/maestro/internal/blueprint"
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
				assert.DeepEqual(t, bp.Name, res.Name, "blueprint name")
				assert.DeepEqual(
					t,
					len(bp.Stages),
					len(res.Stages),
					"blueprint stages len")
				for i, s := range bp.Stages {
					assert.DeepEqual(t, s, res.Stages[i], "stage %d equal", i)
				}
				assert.DeepEqual(
					t,
					len(bp.Links),
					len(res.Links),
					"blueprint links len")
				for i, l := range bp.Links {
					assert.DeepEqual(t, l, res.Links[i], "link %d equal", i)
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
				assert.DeepEqual(t, nil, err, "unmarshal error")
				assert.DeepEqual(t, bp.Name, res.Name, "blueprint name")
				assert.DeepEqual(
					t,
					len(bp.Stages),
					len(res.Stages),
					"blueprint stages len")
				for i, s := range bp.Stages {
					assert.DeepEqual(t, s, res.Stages[i], "stage %d equal", i)
				}
				assert.DeepEqual(
					t,
					len(bp.Links),
					len(res.Links),
					"blueprint links len")
				for i, l := range bp.Links {
					assert.DeepEqual(t, l, res.Links[i], "link %d equal", i)
				}
			})
	}
}

func TestUnmarshalBlueprintError(t *testing.T) {
	res, err := UnmarshalBlueprint(nil)
	assert.DeepEqual(t, errors.New("p is nil"), err, "unmarshal error")
	assert.IsNil(t, res, "nil return value")
}
