package get

import (
	"context"
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/cli/client"
	"github.com/DuarteMRAlves/maestro/internal/cli/resources"
	"gotest.tools/v3/assert"
	"testing"
	"time"
)

// assetForNum deterministically creates an asset resource with the given
// number.
func assetForNum(num int) *resources.AssetResource {
	return &resources.AssetResource{
		Name:  assetNameForNum(num),
		Image: assetImageForNum(num),
	}
}

// stageForNum deterministically creates a stage resource with the given number.
// The associated asset name is the one used in assetForNum.
func stageForNum(num int) *resources.StageResource {
	return &resources.StageResource{
		Name:    stageNameForNum(num),
		Asset:   assetNameForNum(num),
		Service: stageServiceForNum(num),
		Method:  stageMethodForNum(num),
	}
}

// linkForNum deterministically creates a link resource with the given number.
// The associated source stage name is the one used in stageForNum with the num
// argument. The associated target stage name is the one used in the stageForNum
// with the num+1 argument.
func linkForNum(num int) *resources.LinkResource {
	return &resources.LinkResource{
		Name:        linkNameForNum(num),
		SourceStage: linkSourceStageForNum(num),
		SourceField: linkSourceFieldForNum(num),
		TargetStage: linkTargetStageForNum(num),
		TargetField: linkTargetFieldForNum(num),
	}
}

// assetNameForNum deterministically creates an asset name for a given number.
func assetNameForNum(num int) string {
	return fmt.Sprintf("asset-%v", num)
}

// assetImageForNum deterministically creates an image for a given number.
func assetImageForNum(num int) string {
	name := assetNameForNum(num)
	return fmt.Sprintf("image-%v", name)
}

// stageNameForNum deterministically creates a stage name for a given number.
func stageNameForNum(num int) string {
	return fmt.Sprintf("stage-%v", num)
}

// stageServiceForNum deterministically creates a stage service for a given
// number.
func stageServiceForNum(num int) string {
	return fmt.Sprintf("service-%v", num)
}

// stageMethodForNum deterministically creates a stage method for a given
// number.
func stageMethodForNum(num int) string {
	return fmt.Sprintf("method-%v", num)
}

// linkNameForNum deterministically creates a link name for a given number.
func linkNameForNum(num int) string {
	return fmt.Sprintf("link-%v", num)
}

// linkSourceStageForNum deterministically creates a link source stage for a
// given number.
func linkSourceStageForNum(num int) string {
	return stageNameForNum(num)
}

// linkSourceFieldForNum deterministically creates a link source field for a
// given number.
func linkSourceFieldForNum(num int) string {
	return fmt.Sprintf("source-field-%v", num)
}

// linkTargetStageForNum deterministically creates a link target stage for a
// given number.
func linkTargetStageForNum(num int) string {
	return stageNameForNum(num + 1)
}

// linkTargetFieldForNum deterministically creates a link target field for a
// given number.
func linkTargetFieldForNum(num int) string {
	return fmt.Sprintf("target-field-%v", num)
}

// populateAssets creates the assets in the server, asserting any occurred
// errors.
func populateAssets(
	t *testing.T,
	assets []*resources.AssetResource,
	addr string,
) {
	// Create asset before executing command
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second)
	defer cancel()

	for _, a := range assets {
		err := client.CreateAsset(ctx, a, addr)
		assert.NilError(t, err, "populate with assets")
	}
}

// populateStages creates the stages in the server, asserting any occurred
// errors.
func populateStages(
	t *testing.T,
	stages []*resources.StageResource,
	addr string,
) {
	// Create asset before executing command
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second)
	defer cancel()

	for _, s := range stages {
		err := client.CreateStage(ctx, s, addr)
		assert.NilError(t, err, "populate with stages")
	}
}

// populateLinks creates the links in the server, asserting any occurred errors.
func populateLinks(t *testing.T, links []*resources.LinkResource, addr string) {
	// Create asset before executing command
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second)
	defer cancel()

	for _, l := range links {
		err := client.CreateLink(ctx, l, addr)
		assert.NilError(t, err, "populate with links")
	}
}
