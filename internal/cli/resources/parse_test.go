package resources

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"gotest.tools/v3/assert"
	"testing"
)

const (
	assetName  = "assetName"
	assetImage = "assetImage"

	stageName    = "stageName"
	stageAsset   = "stageAsset"
	stageService = "stageService"
	stageMethod  = "stageMethod"

	linkName        = "linkName"
	linkSourceStage = "linkSourceStage"
	linkSourceField = "linkSourceField"
	linkTargetStage = "linkTargetStage"
	linkTargetField = "linkTargetField"
)

var fileContent = fmt.Sprintf(
	`
---
kind: asset
spec:
  name: %v
  image: %v
---
kind: stage
spec:
  name: %v
  asset: %v
  service: %v
  method: %v
---
kind: link
spec:
  name: %v
  source_stage: %v
  source_field: %v
  target_stage: %v
  target_field: %v`,
	assetName,
	assetImage,
	stageName,
	stageAsset,
	stageService,
	stageMethod,
	linkName,
	linkSourceStage,
	linkSourceField,
	linkTargetStage,
	linkTargetField)

var expected = []*Resource{
	{
		Kind: assetKind,
		Spec: &AssetSpec{Name: assetName, Image: assetImage},
	},
	{
		Kind: stageKind,
		Spec: &StageSpec{
			Name:    stageName,
			Asset:   stageAsset,
			Service: stageService,
			Method:  stageMethod,
		},
	},
	{
		Kind: linkKind,
		Spec: &LinkSpec{
			Name:        linkName,
			SourceStage: linkSourceStage,
			SourceField: linkSourceField,
			TargetStage: linkTargetStage,
			TargetField: linkTargetField,
		},
	},
}

func TestParseAllTypes(t *testing.T) {
	data := []byte(fileContent)

	resources, err := ParseBytes(data)
	assert.NilError(t, err, "err is not nil: %v", err)
	assert.DeepEqual(t, expected, resources)
}

func TestParseInvalidArguments(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		expectedMsg string
	}{
		{
			name: "missing kind field",
			data: []byte(`
---
spec:
  name: asset-name
  image: image-name
`),
			expectedMsg: "kind not specified",
		},
		{
			name: "empty kind field",
			data: []byte(`
---
kind:
spec:
  name: asset-name
  image: image-name
`),
			expectedMsg: "kind not specified",
		},
		{
			name: "unknown kind",
			data: []byte(`
---
kind: unknown
spec:
  name: asset-name
  image: image-name
`),
			expectedMsg: "unknown kind: 'unknown'",
		},
		{
			name: "empty spec",
			data: []byte(`
---
kind: stage
spec:
`),
			expectedMsg: "empty spec",
		},
		{
			name: "missing asset name",
			data: []byte(`
---
kind: asset
spec:
  image: image-name
`),
			expectedMsg: "missing required field: 'name'",
		},
		{
			name: "missing stage name",
			data: []byte(`
---
kind: stage
spec:
  asset: stage-asset
`),
			expectedMsg: "missing required field: 'name'",
		},
		{
			name: "missing link name",
			data: []byte(`
---
kind: link
spec:
  source_stage: source-stage
  target_stage: target-stage
`),
			expectedMsg: "missing required field: 'name'",
		},
		{
			name: "missing link source stage",
			data: []byte(`
---
kind: link
spec:
  name: link-name
  target_stage: target-stage
`),
			expectedMsg: "missing required field: 'source_stage'",
		},
		{
			name: "missing link target stage",
			data: []byte(`
---
kind: link
spec:
  name: link-name
  source_stage: source-stage
`),
			expectedMsg: "missing required field: 'target_stage'",
		},
	}
	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				resources, err := ParseBytes(test.data)

				assert.Assert(t, errdefs.IsInvalidArgument(err), "err type")
				assert.ErrorContains(t, err, test.expectedMsg)
				assert.Assert(t, resources == nil, "resources not nil")
			})
	}
}
