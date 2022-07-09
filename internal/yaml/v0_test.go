package yaml

import (
	"errors"
	"reflect"
	"testing"

	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/google/go-cmp/cmp"
)

func TestReadV0(t *testing.T) {
	file := "../../test/data/unit/read/v0/correct.yml"
	resources, err := ReadV0(file)
	if err != nil {
		t.Fatalf("read error: %s", err)
	}

	expected := &api.Pipeline{
		Name: "v0-pipeline",
		Mode: api.OnlineExecution,
		Stages: []*api.Stage{
			{
				Name: "stage-1",
				MethodContext: api.MethodContext{
					Address: "host-1:1",
				},
			},
			{
				Name: "stage-2",
				MethodContext: api.MethodContext{
					Address: "host-2:2",
					Service: "Service2",
				},
			},
			{
				Name: "stage-3",
				MethodContext: api.MethodContext{
					Address: "host-3:3",
					Method:  "Method3",
				},
			},
			{
				Name: "stage-4",
				MethodContext: api.MethodContext{
					Address: "host-4:4",
					Service: "Service4",
					Method:  "Method4",
				},
			},
		},
		Links: []*api.Link{
			{
				Name:        "v0-link-stage-1-to-stage-2",
				SourceStage: "stage-1",
				TargetStage: "stage-2",
			},
			{
				Name:        "v0-link-stage-2-to-stage-3",
				SourceStage: "stage-2",
				SourceField: "Field2",
				TargetStage: "stage-3",
			},
			{
				Name:        "v0-link-stage-3-to-stage-4",
				SourceStage: "stage-3",
				TargetStage: "stage-4",
				TargetField: "Field4",
			},
			{
				Name:        "v0-link-stage-4-to-stage-1",
				SourceStage: "stage-4",
				SourceField: "Field4",
				TargetStage: "stage-1",
				TargetField: "Field1",
			},
		},
	}
	if diff := cmp.Diff(expected, resources); diff != "" {
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
			pipeline, err := ReadV0(tc.file)
			if err == nil {
				t.Fatalf("expected error but got nil")
			}
			if pipeline != nil {
				t.Fatalf("expected nil pipeline")
			}
			tc.verifyErr(t, err)
		})
	}
}
