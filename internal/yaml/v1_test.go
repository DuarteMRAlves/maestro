package yaml

import (
	"errors"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/google/go-cmp/cmp"
)

func TestReadV1(t *testing.T) {
	tests := map[string]struct {
		files    []string
		expected []*api.Pipeline
	}{
		"single file": {
			files: []string{"../../test/data/unit/read/v1/read_single_file.yml"},
			expected: []*api.Pipeline{
				{
					Name: "pipeline-2",
					Mode: api.OfflineExecution,
				},
				{
					Name: "pipeline-1",
					Mode: api.OnlineExecution,
					Stages: []*api.Stage{
						{
							Name: "stage-1",
							MethodContext: api.MethodContext{
								Address: "address-1",
								Service: "Service1",
								Method:  "Method1",
							},
						},
						{
							Name: "stage-2",
							MethodContext: api.MethodContext{
								Address: "address-2",
								Service: "Service2",
							},
						},
						{
							Name: "stage-3",
							MethodContext: api.MethodContext{
								Address: "address-3",
								Method:  "Method3",
							},
						},
					},
					Links: []*api.Link{
						{
							Name:             "link-stage-2-stage-1",
							SourceStage:      "stage-2",
							TargetStage:      "stage-1",
							NumEmptyMessages: 2,
						},
						{
							Name:        "link-stage-1-stage-2",
							SourceStage: "stage-1",
							SourceField: "Field1",
							TargetStage: "stage-2",
							TargetField: "Field2",
						},
					},
				},
			},
		},
		"multiple files": {
			files: []string{
				"../../test/data/unit/read/v1/multi_file1.yml",
				"../../test/data/unit/read/v1/multi_file2.yml",
				"../../test/data/unit/read/v1/multi_file3.yml",
			},
			expected: []*api.Pipeline{
				{
					Name: "pipeline-3",
					Mode: api.OfflineExecution,
					Stages: []*api.Stage{
						{
							Name: "stage-5",
							MethodContext: api.MethodContext{
								Address: "address-5",
								Method:  "Method5",
							},
						},
						{
							Name: "stage-6",
							MethodContext: api.MethodContext{
								Address: "address-6",
								Service: "Service6",
								Method:  "Method6",
							},
						},
					},
					Links: []*api.Link{
						{
							Name:        "link-stage-5-stage-6",
							SourceStage: "stage-5",
							TargetStage: "stage-6",
							TargetField: "Field1",
						},
					},
				},
				{
					Name: "pipeline-4",
					Mode: api.OfflineExecution,
					Stages: []*api.Stage{
						{
							Name: "stage-4",
							MethodContext: api.MethodContext{
								Address: "address-4",
							},
						},
						{
							Name: "stage-7",
							MethodContext: api.MethodContext{
								Address: "address-7",
								Service: "Service7",
							},
						},
					},
					Links: []*api.Link{
						{
							Name:        "link-stage-4-stage-5",
							SourceStage: "stage-4",
							TargetStage: "stage-5",
						},
						{
							Name:        "link-stage-4-stage-6",
							SourceStage: "stage-4",
							TargetStage: "stage-6",
							TargetField: "Field2",
						},
					},
				},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			resources, err := ReadV1(tc.files...)
			if err != nil {
				t.Fatalf("read error: %s", err)
			}
			if diff := cmp.Diff(tc.expected, resources); diff != "" {
				t.Fatalf("read resources mismatch:\n%s", diff)
			}
		})
	}
}

func TestReadV1_Err(t *testing.T) {
	tests := map[string]struct {
		files     []string
		verifyErr func(t *testing.T, err error)
	}{
		"missing kind": {
			files: []string{"../../test/data/unit/read/v1/err_missing_kind.yml"},
			verifyErr: func(t *testing.T, err error) {
				expErr := ErrMissingKind
				if !errors.Is(err, expErr) {
					t.Fatalf("Wrong error: expected '%s', got '%s'", expErr, err)
				}
			},
		},
		"empty spec": {
			files: []string{"../../test/data/unit/read/v1/err_empty_spec.yml"},
			verifyErr: func(t *testing.T, err error) {
				expErr := ErrEmptySpec
				if !errors.Is(err, expErr) {
					t.Fatalf("Wrong error: expected '%s', got '%s'", expErr, err)
				}
			},
		},
		"unknown kind": {
			files: []string{"../../test/data/unit/read/v1/err_unknown_kind.yml"},
			verifyErr: func(t *testing.T, err error) {
				var actual *unknownKind
				if !errors.As(err, &actual) {
					format := "Wrong error type: expected *unknownKind, got %s"
					t.Fatalf(format, reflect.TypeOf(err))
				}
				expected := &unknownKind{Kind: "unknown-kind"}
				if diff := cmp.Diff(expected, actual); diff != "" {
					t.Fatalf("error mismatch:\n%s", diff)
				}
			},
		},
		"missing required field": {
			files: []string{"../../test/data/unit/read/v1/err_missing_req_field.yml"},
			verifyErr: func(t *testing.T, err error) {
				var actual *missingRequiredField
				if !errors.As(err, &actual) {
					format := "Wrong error type: expected *missingRequiredField, got %s"
					t.Fatalf(format, reflect.TypeOf(err))
				}
				expected := &missingRequiredField{Field: "address"}
				if diff := cmp.Diff(expected, actual); diff != "" {
					t.Fatalf("error mismatch:\n%s", diff)
				}
			},
		},
		"unknown fields": {
			files: []string{"../../test/data/unit/read/v1/err_unknown_fields.yml"},
			verifyErr: func(t *testing.T, err error) {
				var actual *unknownFields
				if !errors.As(err, &actual) {
					format := "Wrong error type: expected *unknownFields, got %s"
					t.Fatalf(format, reflect.TypeOf(err))
				}
				expected := &unknownFields{Fields: []string{"unknown_1", "unknown_2"}}
				if diff := cmp.Diff(expected, actual); diff != "" {
					t.Fatalf("error mismatch:\n%s", diff)
				}
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			pipelines, err := ReadV1(tc.files...)
			if err == nil {
				t.Fatalf("expected error but got nil")
			}
			if pipelines != nil {
				t.Fatalf("resources not nil")
			}
			tc.verifyErr(t, err)
		})
	}
}

func TestWriteV1(t *testing.T) {
	pipeline := api.Pipeline{
		Name: "pipeline-1",
		Mode: api.OnlineExecution,
		Stages: []*api.Stage{
			{
				Name: "stage-1",
				MethodContext: api.MethodContext{
					Address: "address-1",
					Service: "Service1",
					Method:  "Method1",
				},
			},
			{
				Name: "stage-2",
				MethodContext: api.MethodContext{
					Address: "address-2",
					Service: "Service2",
				},
			},
			{
				Name: "stage-3",
				MethodContext: api.MethodContext{
					Address: "address-3",
					Method:  "Method3",
				},
			},
		},
		Links: []*api.Link{
			{
				Name:             "link-stage-2-stage-1",
				SourceStage:      "stage-2",
				TargetStage:      "stage-1",
				NumEmptyMessages: 2,
			},
			{
				Name:        "link-stage-1-stage-2",
				SourceStage: "stage-1",
				SourceField: "Field1",
				TargetStage: "stage-2",
				TargetField: "Field2",
			},
		},
	}
	tempDir := t.TempDir()
	outFile := tempDir + "/to_v1.yml"
	err := WriteV1(&pipeline, outFile, 0777)
	if err != nil {
		t.Fatalf("write v1: %s", err)
	}
	writeData, err := ioutil.ReadFile(outFile)
	if err != nil {
		t.Fatalf("read new file: %s", err)
	}
	writeContent := string(writeData)

	expFile := "../../test/data/unit/read/v1/write_single_file.yml"
	expData, err := ioutil.ReadFile(expFile)
	if err != nil {
		t.Fatalf("read v1: %s", err)
	}
	expContent := string(expData)

	if diff := cmp.Diff(expContent, writeContent); diff != "" {
		t.Fatalf("content mismatch:\n%s", diff)
	}
}
