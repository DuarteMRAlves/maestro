package parse

import (
	"github.com/DuarteMRAlves/maestro/internal"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestFromV0(t *testing.T) {
	file := "../../test/data/unit/parse/v0/correct.yml"
	resources, err := FromV0(file)
	if err != nil {
		t.Fatalf("parse error: %s", err)
	}

	expected := ResourceSet{
		Orchestrations: []Orchestration{createOrchestration(t, "v0-orchestration")},
		Stages: []Stage{
			createStage(t, "stage-1", "host-1:1", "", "", "v0-orchestration"),
			createStage(t, "stage-2", "host-2:2", "Service2", "", "v0-orchestration"),
			createStage(t, "stage-3", "host-3:3", "", "Method3", "v0-orchestration"),
			createStage(t, "stage-4", "host-4:4", "Service4", "Method4", "v0-orchestration"),
		},
		Links: []Link{
			createLink(
				t,
				"v0-link-stage-1-to-stage-2",
				"stage-1",
				"",
				"stage-2",
				"",
				"v0-orchestration",
			),
			createLink(
				t,
				"v0-link-stage-2-to-stage-3",
				"stage-2",
				"Field2",
				"stage-3",
				"",
				"v0-orchestration",
			),
			createLink(
				t,
				"v0-link-stage-3-to-stage-4",
				"stage-3",
				"",
				"stage-4",
				"Field4",
				"v0-orchestration",
			),
			createLink(
				t,
				"v0-link-stage-4-to-stage-1",
				"stage-4",
				"Field4",
				"stage-1",
				"Field1",
				"v0-orchestration",
			),
		},
		Assets: nil,
	}
	cmpOpts := cmp.AllowUnexported(
		internal.AssetName{},
		internal.Image{},
		internal.StageName{},
		internal.Address{},
		internal.Service{},
		internal.Method{},
		internal.LinkName{},
		internal.MessageField{},
		internal.OrchestrationName{},
	)
	if diff := cmp.Diff(expected, resources, cmpOpts); diff != "" {
		t.Fatalf("parsed resources mismatch:\n%s", diff)
	}
}
