package parse

import (
	"errors"
	"github.com/DuarteMRAlves/maestro/internal"
	"github.com/google/go-cmp/cmp"
	"io/ioutil"
	"reflect"
	"testing"
)

func TestFromV1(t *testing.T) {
	tests := map[string]struct {
		files    []string
		expected ResourceSet
	}{
		"single file": {
			files: []string{"../../test/data/unit/parse/v1/single_file.yml"},
			expected: ResourceSet{
				Orchestrations: []Orchestration{
					createOrchestration(t, "orchestration-2"),
					createOrchestration(t, "orchestration-1"),
				},
				Stages: []Stage{
					createStage(
						t, "stage-1", "address-1", "Service1", "Method1", "orchestration-1",
					),
					createStage(
						t, "stage-2", "address-2", "Service2", "", "orchestration-1",
					),
					createStage(
						t, "stage-3", "address-3", "", "Method3", "orchestration-1",
					),
				},
				Links: []Link{
					createLink(
						t,
						"link-stage-2-stage-1",
						"stage-2",
						"",
						"stage-1",
						"",
						"orchestration-1",
					),
					createLink(
						t,
						"link-stage-1-stage-2",
						"stage-1",
						"Field1",
						"stage-2",
						"Field2",
						"orchestration-1",
					),
				},
				Assets: []Asset{
					createAsset(t, "asset-1", "image-1"),
					createAsset(t, "asset-2", ""),
				},
			},
		},
		"multiple files": {
			files: []string{
				"../../test/data/unit/parse/v1/multi_file1.yml",
				"../../test/data/unit/parse/v1/multi_file2.yml",
				"../../test/data/unit/parse/v1/multi_file3.yml",
				"../../test/data/unit/parse/v1/multi_file4.yml",
			},
			expected: ResourceSet{
				Orchestrations: []Orchestration{
					createOrchestration(t, "orchestration-3"),
					createOrchestration(t, "orchestration-4"),
				},
				Stages: []Stage{
					createStage(
						t, "stage-4", "address-4", "", "", "orchestration-4",
					),
					createStage(
						t, "stage-5", "address-5", "", "Method5", "orchestration-3",
					),
					createStage(
						t, "stage-6", "address-6", "Service6", "Method6", "orchestration-3",
					),
					createStage(
						t, "stage-7", "address-7", "Service7", "", "orchestration-4",
					),
				},
				Links: []Link{
					createLink(
						t,
						"link-stage-4-stage-5",
						"stage-4",
						"",
						"stage-5",
						"",
						"orchestration-4",
					),
					createLink(
						t,
						"link-stage-5-stage-6",
						"stage-5",
						"",
						"stage-6",
						"Field1",
						"orchestration-3",
					),
					createLink(
						t,
						"link-stage-4-stage-6",
						"stage-4",
						"",
						"stage-6",
						"Field2",
						"orchestration-4",
					),
				},
				Assets: []Asset{
					createAsset(t, "asset-4", "image-4"),
					createAsset(t, "asset-5", ""),
					createAsset(t, "asset-6", "image-6"),
				},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			resources, err := FromV1(tc.files...)
			if err != nil {
				t.Fatalf("parse error: %s", err)
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
			if diff := cmp.Diff(tc.expected, resources, cmpOpts); diff != "" {
				t.Fatalf("parsed resources mismatch:\n%s", diff)
			}
		})
	}
}

func TestFromV1_Err(t *testing.T) {
	tests := map[string]struct {
		files     []string
		verifyErr func(t *testing.T, err error)
	}{
		"missing kind": {
			files: []string{"../../test/data/unit/parse/v1/err_missing_kind.yml"},
			verifyErr: func(t *testing.T, err error) {
				expErr := MissingKind
				if !errors.Is(err, expErr) {
					t.Fatalf("Wrong error: expected '%s', got '%s'", expErr, err)
				}
			},
		},
		"empty spec": {
			files: []string{"../../test/data/unit/parse/v1/err_empty_spec.yml"},
			verifyErr: func(t *testing.T, err error) {
				expErr := EmptySpec
				if !errors.Is(err, expErr) {
					t.Fatalf("Wrong error: expected '%s', got '%s'", expErr, err)
				}
			},
		},
		"unknown kind": {
			files: []string{"../../test/data/unit/parse/v1/err_unknown_kind.yml"},
			verifyErr: func(t *testing.T, err error) {
				var actual *UnknownKind
				if !errors.As(err, &actual) {
					format := "Wrong error type: expected *UnknownKind, got %s"
					t.Fatalf(format, reflect.TypeOf(err))
				}
				expected := &UnknownKind{Kind: "unknown-kind"}
				if diff := cmp.Diff(expected, actual); diff != "" {
					t.Fatalf("error mismatch:\n%s", diff)
				}
			},
		},
		"missing required field": {
			files: []string{"../../test/data/unit/parse/v1/err_missing_req_field.yml"},
			verifyErr: func(t *testing.T, err error) {
				var actual *MissingRequiredField
				if !errors.As(err, &actual) {
					format := "Wrong error type: expected *MissingRequiredField, got %s"
					t.Fatalf(format, reflect.TypeOf(err))
				}
				expected := &MissingRequiredField{Field: "address"}
				if diff := cmp.Diff(expected, actual); diff != "" {
					t.Fatalf("error mismatch:\n%s", diff)
				}
			},
		},
		"unknown fields": {
			files: []string{"../../test/data/unit/parse/v1/err_unknown_fields.yml"},
			verifyErr: func(t *testing.T, err error) {
				var actual *UnknownFields
				if !errors.As(err, &actual) {
					format := "Wrong error type: expected *UnknownFields, got %s"
					t.Fatalf(format, reflect.TypeOf(err))
				}
				expected := &UnknownFields{Fields: []string{"unknown_1", "unknown_2"}}
				if diff := cmp.Diff(expected, actual); diff != "" {
					t.Fatalf("error mismatch:\n%s", diff)
				}
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var emptyResources ResourceSet
			resources, err := FromV1(tc.files...)
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
	resources.Orchestrations = []Orchestration{
		createOrchestration(t, "orchestration-2"),
		createOrchestration(t, "orchestration-1"),
	}
	resources.Stages = []Stage{
		createStage(t, "stage-1", "address-1", "Service1", "Method1", "orchestration-1"),
		createStage(t, "stage-2", "address-2", "Service2", "", "orchestration-1"),
		createStage(t, "stage-3", "address-3", "", "Method3", "orchestration-1"),
	}
	resources.Links = []Link{
		createLink(
			t, "link-stage-2-stage-1", "stage-2", "", "stage-1", "", "orchestration-1",
		),
		createLink(
			t, "link-stage-1-stage-2", "stage-1", "Field1", "stage-2", "Field2", "orchestration-1",
		),
	}
	resources.Assets = []Asset{
		createAsset(t, "asset-1", "image-1"),
		createAsset(t, "asset-2", ""),
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

	expFile := "../../test/data/unit/parse/v1/write_single_file.yml"
	expData, err := ioutil.ReadFile(expFile)
	expContent := string(expData)

	if diff := cmp.Diff(expContent, writeContent); diff != "" {
		t.Fatalf("content mismatch:\n%s", diff)
	}
}

func createOrchestration(t *testing.T, name string) Orchestration {
	orchName, err := internal.NewOrchestrationName(name)
	if err != nil {
		t.Fatalf("create orchestration name %s: %s", name, err)
	}
	return Orchestration{Name: orchName}
}

func createStage(t *testing.T, name, addr, serv, meth, orch string) Stage {
	stageName, err := internal.NewStageName(name)
	if err != nil {
		t.Fatalf("create stage name %s: %s", name, err)
	}
	methodCtx := MethodContext{
		Address: internal.NewAddress(addr),
		Service: internal.NewService(serv),
		Method:  internal.NewMethod(meth),
	}
	orchName, err := internal.NewOrchestrationName(orch)
	if err != nil {
		t.Fatalf("create orchestration name %s: %s", orch, err)
	}
	return Stage{Name: stageName, Method: methodCtx, Orchestration: orchName}
}

func createLink(
	t *testing.T, name, srcStage, srcField, tgtStage, tgtField, orch string,
) Link {
	linkName, err := internal.NewLinkName(name)
	if err != nil {
		t.Fatalf("create link name %s: %s", name, err)
	}
	srcName, err := internal.NewStageName(srcStage)
	if err != nil {
		t.Fatalf("create stage name %s: %s", name, err)
	}
	tgtName, err := internal.NewStageName(tgtStage)
	if err != nil {
		t.Fatalf("create stage name %s: %s", name, err)
	}
	srcEndPt := LinkEndpoint{
		Stage: srcName, Field: internal.NewMessageField(srcField),
	}
	tgtEndPt := LinkEndpoint{
		Stage: tgtName, Field: internal.NewMessageField(tgtField),
	}
	orchName, err := internal.NewOrchestrationName(orch)
	if err != nil {
		t.Fatalf("create orchestration name %s: %s", orch, err)
	}
	return Link{
		Name: linkName, Source: srcEndPt, Target: tgtEndPt, Orchestration: orchName,
	}
}

func createAsset(t *testing.T, name, img string) Asset {
	assetName, err := internal.NewAssetName(name)
	if err != nil {
		t.Fatalf("create asset name %s: %s", name, err)
	}
	image := internal.NewImage(img)
	return Asset{Name: assetName, Image: image}
}
