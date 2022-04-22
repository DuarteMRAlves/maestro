package yaml

import (
	"errors"
	"github.com/DuarteMRAlves/maestro/internal"
	"github.com/google/go-cmp/cmp"
	"io/ioutil"
	"reflect"
	"testing"
)

func TestReadV1(t *testing.T) {
	tests := map[string]struct {
		files    []string
		expected ResourceSet
	}{
		"single file": {
			files: []string{"../../test/data/unit/read/v1/read_single_file.yml"},
			expected: ResourceSet{
				Pipelines: []Pipeline{
					{Name: newV1PipelineName(t, "pipeline-2"), Mode: internal.OfflineExecution},
					{Name: newV1PipelineName(t, "pipeline-1"), Mode: internal.OnlineExecution},
				},
				Stages: []Stage{
					{
						Name: newV1StageName(t, "stage-1"),
						Method: MethodContext{
							Address: internal.NewAddress("address-1"),
							Service: internal.NewService("Service1"),
							Method:  internal.NewMethod("Method1"),
						},
						Pipeline: newV1PipelineName(t, "pipeline-1"),
					},
					{
						Name: newV1StageName(t, "stage-2"),
						Method: MethodContext{
							Address: internal.NewAddress("address-2"),
							Service: internal.NewService("Service2"),
						},
						Pipeline: newV1PipelineName(t, "pipeline-1"),
					},
					{
						Name: newV1StageName(t, "stage-3"),
						Method: MethodContext{
							Address: internal.NewAddress("address-3"),
							Method:  internal.NewMethod("Method3"),
						},
						Pipeline: newV1PipelineName(t, "pipeline-1"),
					},
				},
				Links: []Link{
					{
						Name:     newV1LinkName(t, "link-stage-2-stage-1"),
						Source:   LinkEndpoint{Stage: newV1StageName(t, "stage-2")},
						Target:   LinkEndpoint{Stage: newV1StageName(t, "stage-1")},
						Pipeline: newV1PipelineName(t, "pipeline-1"),
					},
					{
						Name: newV1LinkName(t, "link-stage-1-stage-2"),
						Source: LinkEndpoint{
							Stage: newV1StageName(t, "stage-1"),
							Field: internal.NewMessageField("Field1"),
						},
						Target: LinkEndpoint{
							Stage: newV1StageName(t, "stage-2"),
							Field: internal.NewMessageField("Field2"),
						},
						Pipeline: newV1PipelineName(t, "pipeline-1"),
					},
				},
				Assets: []Asset{
					{Name: newV1AssetName(t, "asset-1"), Image: internal.NewImage("image-1")},
					{Name: newV1AssetName(t, "asset-2"), Image: internal.NewImage("")},
				},
			},
		},
		"multiple files": {
			files: []string{
				"../../test/data/unit/read/v1/multi_file1.yml",
				"../../test/data/unit/read/v1/multi_file2.yml",
				"../../test/data/unit/read/v1/multi_file3.yml",
				"../../test/data/unit/read/v1/multi_file4.yml",
			},
			expected: ResourceSet{
				Pipelines: []Pipeline{
					{Name: newV1PipelineName(t, "pipeline-3"), Mode: internal.OfflineExecution},
					{Name: newV1PipelineName(t, "pipeline-4"), Mode: internal.OfflineExecution},
				},
				Stages: []Stage{
					{
						Name: newV1StageName(t, "stage-4"),
						Method: MethodContext{
							Address: internal.NewAddress("address-4"),
						},
						Pipeline: newV1PipelineName(t, "pipeline-4"),
					},
					{
						Name: newV1StageName(t, "stage-5"),
						Method: MethodContext{
							Address: internal.NewAddress("address-5"),
							Method:  internal.NewMethod("Method5"),
						},
						Pipeline: newV1PipelineName(t, "pipeline-3"),
					},
					{
						Name: newV1StageName(t, "stage-6"),
						Method: MethodContext{
							Address: internal.NewAddress("address-6"),
							Service: internal.NewService("Service6"),
							Method:  internal.NewMethod("Method6"),
						},
						Pipeline: newV1PipelineName(t, "pipeline-3"),
					},
					{
						Name: newV1StageName(t, "stage-7"),
						Method: MethodContext{
							Address: internal.NewAddress("address-7"),
							Service: internal.NewService("Service7"),
						},
						Pipeline: newV1PipelineName(t, "pipeline-4"),
					},
				},
				Links: []Link{
					{
						Name: newV1LinkName(t, "link-stage-4-stage-5"),
						Source: LinkEndpoint{
							Stage: newV1StageName(t, "stage-4"),
						},
						Target: LinkEndpoint{
							Stage: newV1StageName(t, "stage-5"),
						},
						Pipeline: newV1PipelineName(t, "pipeline-4"),
					},
					{
						Name: newV1LinkName(t, "link-stage-5-stage-6"),
						Source: LinkEndpoint{
							Stage: newV1StageName(t, "stage-5"),
						},
						Target: LinkEndpoint{
							Stage: newV1StageName(t, "stage-6"),
							Field: internal.NewMessageField("Field1"),
						},
						Pipeline: newV1PipelineName(t, "pipeline-3"),
					},
					{
						Name: newV1LinkName(t, "link-stage-4-stage-6"),
						Source: LinkEndpoint{
							Stage: newV1StageName(t, "stage-4"),
						},
						Target: LinkEndpoint{
							Stage: newV1StageName(t, "stage-6"),
							Field: internal.NewMessageField("Field2"),
						},
						Pipeline: newV1PipelineName(t, "pipeline-4"),
					},
				},
				Assets: []Asset{
					{Name: newV1AssetName(t, "asset-4"), Image: internal.NewImage("image-4")},
					{Name: newV1AssetName(t, "asset-5"), Image: internal.NewImage("")},
					{Name: newV1AssetName(t, "asset-6"), Image: internal.NewImage("image-6")},
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
			if diff := cmp.Diff(tc.expected, resources, cmpOpts); diff != "" {
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
				expErr := MissingKind
				if !errors.Is(err, expErr) {
					t.Fatalf("Wrong error: expected '%s', got '%s'", expErr, err)
				}
			},
		},
		"empty spec": {
			files: []string{"../../test/data/unit/read/v1/err_empty_spec.yml"},
			verifyErr: func(t *testing.T, err error) {
				expErr := EmptySpec
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
			var emptyResources ResourceSet
			resources, err := ReadV1(tc.files...)
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

func TestWriteV1(t *testing.T) {
	var resources ResourceSet
	resources.Pipelines = []Pipeline{
		{Name: newV1PipelineName(t, "pipeline-2"), Mode: internal.OnlineExecution},
		{Name: newV1PipelineName(t, "pipeline-1"), Mode: internal.OfflineExecution},
	}
	resources.Stages = []Stage{
		{
			Name: newV1StageName(t, "stage-1"),
			Method: MethodContext{
				Address: internal.NewAddress("address-1"),
				Service: internal.NewService("Service1"),
				Method:  internal.NewMethod("Method1"),
			},
			Pipeline: newV1PipelineName(t, "pipeline-1"),
		},
		{
			Name: newV1StageName(t, "stage-2"),
			Method: MethodContext{
				Address: internal.NewAddress("address-2"),
				Service: internal.NewService("Service2"),
			},
			Pipeline: newV1PipelineName(t, "pipeline-1"),
		},
		{
			Name: newV1StageName(t, "stage-3"),
			Method: MethodContext{
				Address: internal.NewAddress("address-3"),
				Method:  internal.NewMethod("Method3"),
			},
			Pipeline: newV1PipelineName(t, "pipeline-1"),
		},
	}
	resources.Links = []Link{
		{
			Name:     newV1LinkName(t, "link-stage-2-stage-1"),
			Source:   LinkEndpoint{Stage: newV1StageName(t, "stage-2")},
			Target:   LinkEndpoint{Stage: newV1StageName(t, "stage-1")},
			Pipeline: newV1PipelineName(t, "pipeline-1"),
		},
		{
			Name: newV1LinkName(t, "link-stage-1-stage-2"),
			Source: LinkEndpoint{
				Stage: newV1StageName(t, "stage-1"),
				Field: internal.NewMessageField("Field1"),
			},
			Target: LinkEndpoint{
				Stage: newV1StageName(t, "stage-2"),
				Field: internal.NewMessageField("Field2"),
			},
			Pipeline: newV1PipelineName(t, "pipeline-1"),
		},
	}
	resources.Assets = []Asset{
		{Name: newV1AssetName(t, "asset-1"), Image: internal.NewImage("image-1")},
		{Name: newV1AssetName(t, "asset-2"), Image: internal.NewImage("")},
	}
	tempDir := t.TempDir()
	outFile := tempDir + "/to_v1.yml"
	err := WriteV1(resources, outFile, 777)
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
	expContent := string(expData)

	if diff := cmp.Diff(expContent, writeContent); diff != "" {
		t.Fatalf("content mismatch:\n%s", diff)
	}
}

func newV1LinkName(t *testing.T, name string) internal.LinkName {
	linkName, err := internal.NewLinkName(name)
	if err != nil {
		t.Fatalf("new v1 link name %s: %s", name, err)
	}
	return linkName
}

func newV1StageName(t *testing.T, name string) internal.StageName {
	stageName, err := internal.NewStageName(name)
	if err != nil {
		t.Fatalf("new v1 stage name %s: %s", name, err)
	}
	return stageName
}

func newV1PipelineName(t *testing.T, name string) internal.PipelineName {
	pipelineName, err := internal.NewPipelineName(name)
	if err != nil {
		t.Fatalf("new v1 pipeline name %s: %s", name, err)
	}
	return pipelineName
}

func newV1AssetName(t *testing.T, name string) internal.AssetName {
	assetName, err := internal.NewAssetName(name)
	if err != nil {
		t.Fatalf("new v1 asset name %s: %s", name, err)
	}
	return assetName
}
