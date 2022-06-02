package yaml

import (
	"errors"
	"reflect"
	"testing"

	"github.com/DuarteMRAlves/maestro/internal"
	"github.com/google/go-cmp/cmp"
)

func TestReadV0(t *testing.T) {
	file := "../../test/data/unit/read/v0/correct.yml"
	resources, err := ReadV0(file)
	if err != nil {
		t.Fatalf("read error: %s", err)
	}

	expected := ResourceSet{
		Pipelines: []Pipeline{
			{Name: newV0PipelineName(t, "v0-pipeline"), Mode: internal.OnlineExecution},
		},
		Stages: []Stage{
			{
				Name: newV0StageName(t, "stage-1"),
				Method: MethodContext{
					Address: internal.NewAddress("host-1:1"),
				},
				Pipeline: newV0PipelineName(t, "v0-pipeline"),
			},
			{
				Name: newV0StageName(t, "stage-2"),
				Method: MethodContext{
					Address: internal.NewAddress("host-2:2"),
					Service: internal.NewService("Service2"),
				},
				Pipeline: newV0PipelineName(t, "v0-pipeline"),
			},
			{
				Name: newV0StageName(t, "stage-3"),
				Method: MethodContext{
					Address: internal.NewAddress("host-3:3"),
					Method:  internal.NewMethod("Method3"),
				},
				Pipeline: newV0PipelineName(t, "v0-pipeline"),
			},
			{
				Name: newV0StageName(t, "stage-4"),
				Method: MethodContext{
					Address: internal.NewAddress("host-4:4"),
					Service: internal.NewService("Service4"),
					Method:  internal.NewMethod("Method4"),
				},
				Pipeline: newV0PipelineName(t, "v0-pipeline"),
			},
		},
		Links: []Link{
			{
				Name:     newV0LinkName(t, "v0-link-stage-1-to-stage-2"),
				Source:   LinkEndpoint{Stage: newV0StageName(t, "stage-1")},
				Target:   LinkEndpoint{Stage: newV0StageName(t, "stage-2")},
				Pipeline: newV0PipelineName(t, "v0-pipeline"),
			},
			{
				Name: newV0LinkName(t, "v0-link-stage-2-to-stage-3"),
				Source: LinkEndpoint{
					Stage: newV0StageName(t, "stage-2"),
					Field: internal.NewMessageField("Field2"),
				},
				Target:   LinkEndpoint{Stage: newV0StageName(t, "stage-3")},
				Pipeline: newV0PipelineName(t, "v0-pipeline"),
			},
			{
				Name:   newV0LinkName(t, "v0-link-stage-3-to-stage-4"),
				Source: LinkEndpoint{Stage: newV0StageName(t, "stage-3")},
				Target: LinkEndpoint{
					Stage: newV0StageName(t, "stage-4"),
					Field: internal.NewMessageField("Field4"),
				},
				Pipeline: newV0PipelineName(t, "v0-pipeline"),
			},
			{
				Name: newV0LinkName(t, "v0-link-stage-4-to-stage-1"),
				Source: LinkEndpoint{
					Stage: newV0StageName(t, "stage-4"),
					Field: internal.NewMessageField("Field4"),
				},
				Target: LinkEndpoint{
					Stage: newV0StageName(t, "stage-1"),
					Field: internal.NewMessageField("Field1"),
				},
				Pipeline: newV0PipelineName(t, "v0-pipeline"),
			},
		},
	}
	cmpOpts := cmp.AllowUnexported(
		internal.StageName{},
		internal.Address{},
		internal.Service{},
		internal.Method{},
		internal.LinkName{},
		internal.MessageField{},
		internal.PipelineName{},
		internal.ExecutionMode{},
	)
	if diff := cmp.Diff(expected, resources, cmpOpts); diff != "" {
		t.Fatalf("read resources mismatch:\n%s", diff)
	}
}

func TestReadV0_Err(t *testing.T) {
	tests := map[string]struct {
		file      string
		verifyErr func(t *testing.T, err error)
	}{
		"unknown fields": {
			file: "../../test/data/unit/read/v0/err_unk_file_tag.yml",
			verifyErr: func(t *testing.T, err error) {
				var actual *unknownFields
				if !errors.As(err, &actual) {
					format := "Wrong error type: expected *unknownFields, got %s"
					t.Fatalf(format, reflect.TypeOf(err))
				}
				expected := &unknownFields{
					Fields: []string{"unknown_base", "unknown_link", "unknown_stage"},
				}
				if diff := cmp.Diff(expected, actual); diff != "" {
					t.Fatalf("error mismatch:\n%s", diff)
				}
			},
		},
		"missing required stage field": {
			file: "../../test/data/unit/read/v0/err_missing_req_stage_field.yml",
			verifyErr: func(t *testing.T, err error) {
				var actual *missingRequiredField
				if !errors.As(err, &actual) {
					format := "Wrong error type: expected *missingRequiredField, got %s"
					t.Fatalf(format, reflect.TypeOf(err))
				}
				expected := &missingRequiredField{Field: "host"}
				if diff := cmp.Diff(expected, actual); diff != "" {
					t.Fatalf("error mismatch:\n%s", diff)
				}
			},
		},
		"missing required link field": {
			file: "../../test/data/unit/read/v0/err_missing_req_link_field.yml",
			verifyErr: func(t *testing.T, err error) {
				var actual *missingRequiredField
				if !errors.As(err, &actual) {
					format := "Wrong error type: expected *missingRequiredField, got %s"
					t.Fatalf(format, reflect.TypeOf(err))
				}
				expected := &missingRequiredField{Field: "stage"}
				if diff := cmp.Diff(expected, actual); diff != "" {
					t.Fatalf("error mismatch:\n%s", diff)
				}
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var emptyResources ResourceSet
			resources, err := ReadV0(tc.file)
			if err == nil {
				t.Fatalf("expected error but got nil")
			}
			if diff := cmp.Diff(emptyResources, resources); diff != "" {
				t.Fatalf("resources not empty:\n%s", diff)
			}
			tc.verifyErr(t, err)
		})
	}
}

func newV0LinkName(t *testing.T, name string) internal.LinkName {
	linkName, err := internal.NewLinkName(name)
	if err != nil {
		t.Fatalf("new v0 link name %s: %s", name, err)
	}
	return linkName
}

func newV0StageName(t *testing.T, name string) internal.StageName {
	stageName, err := internal.NewStageName(name)
	if err != nil {
		t.Fatalf("new v0 stage name %s: %s", name, err)
	}
	return stageName
}

func newV0PipelineName(t *testing.T, name string) internal.PipelineName {
	pipelineName, err := internal.NewPipelineName(name)
	if err != nil {
		t.Fatalf("new v0 pipeline name %s: %s", name, err)
	}
	return pipelineName
}
