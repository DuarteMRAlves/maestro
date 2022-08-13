package repr

import (
	"testing"

	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/google/go-cmp/cmp"
)

func TestPipeline(t *testing.T) {
	tests := map[string]struct {
		input    *api.Pipeline
		expected string
	}{
		"empty stages and links": {
			input: &api.Pipeline{
				Name: "Pipeline",
				Mode: api.OnlineExecution,
			},
			expected: `Name: "Pipeline"
ExecutionMode: Online
Stages: []
Links: []`,
		},
		"non empty stages and empty links": {
			input: &api.Pipeline{
				Name: "Pipeline",
				Stages: []*api.Stage{
					{
						Name:    "Stage-1",
						Address: "address-1",
					},
					{
						Name:    "Stage-2",
						Address: "address-2",
						Service: "Service2",
						Method:  "Method2",
					},
				},
			},
			expected: `Name: "Pipeline"
ExecutionMode: Offline
Stages: [
	{
		Name: "Stage-1"
		Address: "address-1"
		Service: ""
		Method: ""
	}
	{
		Name: "Stage-2"
		Address: "address-2"
		Service: "Service2"
		Method: "Method2"
	}
]
Links: []`,
		},
		"non empty stages and links": {
			input: &api.Pipeline{
				Name: "Pipeline",
				Stages: []*api.Stage{
					{
						Name:    "Stage-1",
						Address: "address-1",
					},
					{
						Name:    "Stage-2",
						Address: "address-2",
						Service: "Service2",
						Method:  "Method2",
					},
				},
				Links: []*api.Link{
					{
						Name:        "Stage-1-to-Stage-2",
						SourceStage: "Stage-1",
						TargetStage: "Stage-2",
					},
					{
						Name:             "Stage-2-to-Stage-1",
						SourceStage:      "Stage-2",
						SourceField:      "Field2",
						TargetStage:      "Stage-1",
						TargetField:      "Field1",
						NumEmptyMessages: 1,
					},
				},
			},
			expected: `Name: "Pipeline"
ExecutionMode: Offline
Stages: [
	{
		Name: "Stage-1"
		Address: "address-1"
		Service: ""
		Method: ""
	}
	{
		Name: "Stage-2"
		Address: "address-2"
		Service: "Service2"
		Method: "Method2"
	}
]
Links: [
	{
		Name: "Stage-1-to-Stage-2"
		SourceStage: "Stage-1"
		SourceField: ""
		TargetStage: "Stage-2"
		TargetField: ""
		NumEmptyMessages: 0
	}
	{
		Name: "Stage-2-to-Stage-1"
		SourceStage: "Stage-2"
		SourceField: "Field2"
		TargetStage: "Stage-1"
		TargetField: "Field1"
		NumEmptyMessages: 1
	}
]`,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actual := Pipeline(tc.input)
			if diff := cmp.Diff(tc.expected, actual); diff != "" {
				t.Fatalf("mismatch in result:\n%s", diff)
			}
		})
	}
}
