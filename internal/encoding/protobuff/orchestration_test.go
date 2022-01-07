package protobuff

import (
	"github.com/DuarteMRAlves/maestro/api/pb"
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"gotest.tools/v3/assert"
	"testing"
)

func TestMarshalOrchestration(t *testing.T) {
	const (
		name  apitypes.OrchestrationName = "OrchestrationName"
		phase                            = apitypes.OrchestrationRunning
		link1 apitypes.LinkName          = "Link Name 1"
		link2 apitypes.LinkName          = "Link Name 2"
		link3 apitypes.LinkName          = "Link Name 3"
	)
	tests := []struct {
		name string
		o    *apitypes.Orchestration
	}{
		{
			name: "all fields with non default values",
			o: &apitypes.Orchestration{
				Name:  name,
				Phase: phase,
				Links: []apitypes.LinkName{link1, link2, link3},
			},
		},
		{
			name: "all field with default values",
			o: &apitypes.Orchestration{
				Name:  "",
				Phase: "",
				Links: nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				o := test.o
				res, err := MarshalOrchestration(o)
				assert.NilError(t, err, "marshal error")
				assert.Equal(t, string(o.Name), res.Name, "orchestration name")
				assert.Equal(t, string(o.Phase), res.Phase, "orchestration phase")
				assert.Equal(
					t,
					len(o.Links),
					len(res.Links),
					"orchestration links len")
				for i, l := range o.Links {
					assert.Equal(t, string(l), res.Links[i], "link %d equal", i)
				}
			})
	}
}

func TestUnmarshalOrchestrationCorrect(t *testing.T) {
	const (
		name  = "OrchestrationName"
		phase = string(apitypes.OrchestrationPending)
		link1 = "Link Name 1"
		link2 = "Link Name 2"
		link3 = "Link Name 3"
	)
	tests := []struct {
		name string
		o    *pb.Orchestration
	}{
		{
			name: "all fields with non default values",
			o: &pb.Orchestration{
				Name:  name,
				Phase: phase,
				Links: []string{link1, link2, link3},
			},
		},
		{
			name: "all field with default values",
			o: &pb.Orchestration{
				Name:  "",
				Phase: "",
				Links: nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				o := test.o
				res, err := UnmarshalOrchestration(o)
				assert.Equal(t, nil, err, "unmarshal error")
				assert.Equal(t, o.Name, string(res.Name), "orchestration name")
				assert.Equal(t, o.Phase, string(res.Phase), "orchestration phase")
				assert.Equal(
					t,
					len(o.Links),
					len(res.Links),
					"orchestration links len")
				for i, l := range o.Links {
					assert.Equal(t, l, string(res.Links[i]), "link %d equal", i)
				}
			})
	}
}

func TestMarshalOrchestrationError(t *testing.T) {
	res, err := MarshalOrchestration(nil)
	assert.Assert(t, errdefs.IsInvalidArgument(err), "err type")
	assert.ErrorContains(t, err, "'o' is nil")
	assert.Assert(t, res == nil, "nil return value")
}

func TestUnmarshalOrchestrationNil(t *testing.T) {
	res, err := UnmarshalOrchestration(nil)
	assert.Assert(t, errdefs.IsInvalidArgument(err), "err type")
	assert.ErrorContains(t, err, "'p' is nil")
	assert.Assert(t, res == nil, "nil return value")
}
