package protobuff

import (
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/blueprint"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"gotest.tools/v3/assert"
	"testing"
)

func TestMarshalBlueprint(t *testing.T) {
	tests := []struct {
		name string
		bp   *blueprint.Blueprint
	}{
		{
			name: "all fields with non default values",
			bp: &blueprint.Blueprint{
				Name:  blueprintName,
				Links: []string{bpLink1, bpLink2, bpLink3},
			},
		},
		{
			name: "all field with default values",
			bp: &blueprint.Blueprint{
				Name:  "",
				Links: nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				bp := test.bp
				res := MarshalBlueprint(bp)
				assert.Equal(t, bp.Name, res.Name, "blueprint name")
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
	tests := []struct {
		name string
		bp   *pb.Blueprint
	}{
		{
			name: "all fields with non default values",
			bp: &pb.Blueprint{
				Name:  blueprintName,
				Links: []string{bpLink1, bpLink2, bpLink3},
			},
		},
		{
			name: "all field with default values",
			bp: &pb.Blueprint{
				Name:  "",
				Links: nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				bp := test.bp
				res, err := UnmarshalBlueprint(bp)
				assert.Equal(t, nil, err, "unmarshal error")
				assert.Equal(t, bp.Name, res.Name, "blueprint name")
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
