package server

import (
	"github.com/DuarteMRAlves/maestro/internal/asset"
	"github.com/DuarteMRAlves/maestro/internal/link"
	"github.com/DuarteMRAlves/maestro/internal/stage"
	"github.com/DuarteMRAlves/maestro/internal/testutil"
	"gotest.tools/v3/assert"
	"testing"
)

// assetForNum deterministically creates an asset with the given number.
func assetForNum(num int) *asset.Asset {
	name := testutil.AssetNameForNum(num)
	img := testutil.AssetImageForNum(num)
	return asset.New(name, img)
}

// stageForNum deterministically creates a stage with the given number.
// The associated asset name is the one used in assetForNum.
func stageForNum(num int) *stage.Stage {
	name := testutil.StageNameForNum(num)
	assetName := testutil.AssetNameForNum(num)
	address := testutil.StageAddressForNum(num)
	return stage.New(name, address, assetName, nil)
}

// linkForNum deterministically creates a link with the given number.
// The associated source stage name is the one used in stageForNum with the num
// argument. The associated target stage name is the one used in the stageForNum
// with the num+1 argument.
func linkForNum(num int) *link.Link {
	name := testutil.LinkNameForNum(num)
	sourceStage := testutil.LinkSourceStageForNum(num)
	sourceField := testutil.LinkSourceFieldForNum(num)
	targetStage := testutil.LinkTargetStageForNum(num)
	targetField := testutil.LinkTargetFieldForNum(num)
	return link.New(name, sourceStage, sourceField, targetStage, targetField)
}

// populateAssets creates the assets in the server, asserting any occurred
// errors.
func populateAssets(t *testing.T, s *Server, assets []*asset.Asset) {
	// Bypass CreateAsset verifications
	store := s.assetStore
	for _, a := range assets {
		assert.NilError(t, store.Create(a), "populate with assets")
	}
}

// populateStages creates the stages in the server, asserting any occurred
// errors.
func populateStages(t *testing.T, s *Server, stages []*stage.Stage) {
	// Bypass CreateStage verifications
	store := s.stageStore
	for _, st := range stages {
		assert.NilError(t, store.Create(st), "populate with stages")
	}
}

// populateLinks creates the links in the server, asserting any occurred errors.
func populateLinks(t *testing.T, s *Server, links []*link.Link) {
	// Bypass CreateLink verifications
	store := s.linkStore
	for _, l := range links {
		assert.NilError(t, store.Create(l), "populate with links")
	}
}
