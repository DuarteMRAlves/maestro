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

func TestMarshalBlueprint(t *testing.T) {
	id, err := identifier.Rand(5)
	assert.IsNil(t, err, "create blueprint id")

	tests := []*blueprint.Blueprint{
		{
			Id:     id,
			Name:   blueprintName,
			Stages: []*blueprint.Stage{stage1, stage2, stage3},
			Links:  []*blueprint.Link{link1, link2, link3},
		},
		{
			Id:     identifier.Empty(),
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
				assert.DeepEqual(t, bp.Id.Val, res.Id.Val, "blueprint id")
				assert.DeepEqual(t, bp.Name, res.Name, "blueprint name")
				assert.DeepEqual(
					t,
					len(bp.Stages),
					len(res.Stages),
					"blueprint stages len")
				for i, s := range bp.Stages {
					assertStage(t, s, res.Stages[i])
				}
				assert.DeepEqual(
					t,
					len(bp.Links),
					len(res.Links),
					"blueprint links len")
				for i, l := range bp.Links {
					assertLink(t, l, res.Links[i])
				}
			})
	}
}

func TestUnmarshalBlueprintCorrect(t *testing.T) {
	tests := []*pb.Blueprint{
		{
			Id:     &pb.Id{Val: blueprintId},
			Name:   blueprintName,
			Stages: []*pb.Stage{pbStage1, pbStage2, pbStage3},
			Links:  []*pb.Link{pbLink1, pbLink2, pbLink3},
		},
		{
			Id:     &pb.Id{Val: ""},
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
				assert.DeepEqual(t, bp.Id.Val, res.Id.Val, "blueprint id")
				assert.DeepEqual(t, bp.Name, res.Name, "blueprint name")
				assert.DeepEqual(
					t,
					len(bp.Stages),
					len(res.Stages),
					"blueprint stages len")
				for i, s := range bp.Stages {
					assertPbStage(t, s, res.Stages[i])
				}
				assert.DeepEqual(
					t,
					len(bp.Links),
					len(res.Links),
					"blueprint links len")
				for i, l := range bp.Links {
					assertPbLink(t, l, res.Links[i])
				}
			})
	}
}

func TestUnmarshalBlueprintError(t *testing.T) {
	res, err := UnmarshalBlueprint(nil)
	assert.DeepEqual(t, errors.New("p is nil"), err, "unmarshal error")
	assert.IsNil(t, res, "nil return value")
}
