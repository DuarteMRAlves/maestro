package parse

import (
	"github.com/DuarteMRAlves/maestro/internal"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestFromV1(t *testing.T) {
	tests := map[string]struct {
		files    []string
		expected ResourceSet
	}{
		"single file": {
			files: []string{"../../test/data/unit/parse/v1/resources.yml"},
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
				"../../test/data/unit/parse/v1/orchestrations.yml",
				"../../test/data/unit/parse/v1/stages.yml",
				"../../test/data/unit/parse/v1/links.yml",
				"../../test/data/unit/parse/v1/assets.yml",
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
