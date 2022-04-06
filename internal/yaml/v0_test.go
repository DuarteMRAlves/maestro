package yaml

import (
	"errors"
	"github.com/DuarteMRAlves/maestro/internal"
	"github.com/google/go-cmp/cmp"
	"reflect"
	"testing"
)

func TestReadV0(t *testing.T) {
	file := "../../test/data/unit/read/v0/correct.yml"
	resources, err := ReadV0(file)
	if err != nil {
		t.Fatalf("read error: %s", err)
	}

	expected := ResourceSet{
		Pipelines: []Pipeline{
			createPipeline(t, "v0-pipeline", internal.OnlineExecution),
		},
		Stages: []Stage{
			createStage(t, "stage-1", "host-1:1", "", "", "v0-pipeline"),
			createStage(t, "stage-2", "host-2:2", "Service2", "", "v0-pipeline"),
			createStage(t, "stage-3", "host-3:3", "", "Method3", "v0-pipeline"),
			createStage(t, "stage-4", "host-4:4", "Service4", "Method4", "v0-pipeline"),
		},
		Links: []Link{
			createLink(
				t,
				"v0-link-stage-1-to-stage-2",
				"stage-1",
				"",
				"stage-2",
				"",
				"v0-pipeline",
			),
			createLink(
				t,
				"v0-link-stage-2-to-stage-3",
				"stage-2",
				"Field2",
				"stage-3",
				"",
				"v0-pipeline",
			),
			createLink(
				t,
				"v0-link-stage-3-to-stage-4",
				"stage-3",
				"",
				"stage-4",
				"Field4",
				"v0-pipeline",
			),
			createLink(
				t,
				"v0-link-stage-4-to-stage-1",
				"stage-4",
				"Field4",
				"stage-1",
				"Field1",
				"v0-pipeline",
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
