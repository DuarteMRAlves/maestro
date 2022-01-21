package server

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/asset"
	"github.com/DuarteMRAlves/maestro/internal/link"
	"github.com/DuarteMRAlves/maestro/internal/reflection"
	"github.com/DuarteMRAlves/maestro/internal/stage"
	"github.com/DuarteMRAlves/maestro/internal/testutil"
	mockreflection "github.com/DuarteMRAlves/maestro/internal/testutil/mock/reflection"
	"github.com/jhump/protoreflect/desc"
	"gotest.tools/v3/assert"
	"reflect"
	"testing"
)

// assetForNum deterministically creates an asset with the given number.
func assetForNum(num int) *asset.Asset {
	name := testutil.AssetNameForNum(num)
	img := testutil.AssetImageForNum(num)
	return asset.New(name, img)
}

// mockStage deterministically creates a stage with the given number.
// The associated asset name is the one used in assetForNum.
func mockStage(
	t *testing.T,
	num int,
	req interface{},
	res interface{},
) *stage.Stage {
	reqType := reflect.TypeOf(req)

	reqDesc, err := desc.LoadMessageDescriptorForType(reqType)
	assert.NilError(t, err, fmt.Sprintf("load req desc for stage: %d\n", num))

	reqMsg, err := reflection.NewMessage(reqDesc)
	assert.NilError(t, err, fmt.Sprintf("load req msg for stage: %d\n", num))

	resType := reflect.TypeOf(res)

	resDesc, err := desc.LoadMessageDescriptorForType(resType)
	assert.NilError(t, err, fmt.Sprintf("load res desc for stage: %d\n", num))

	resMsg, err := reflection.NewMessage(resDesc)
	assert.NilError(t, err, fmt.Sprintf("load res desc for stage: %d\n", num))

	return stage.New(
		testutil.StageNameForNum(num),
		testutil.StageAddressForNum(num),
		testutil.AssetNameForNum(num),
		&mockreflection.RPC{
			Name_: testutil.StageRpcForNum(num),
			FQN: fmt.Sprintf(
				"%s/%s",
				testutil.StageServiceForNum(num),
				testutil.StageRpcForNum(num)),
			In:    reqMsg,
			Out:   resMsg,
			Unary: true,
		})
}

// mockLink deterministically creates a link with the given number.
// The associated source stage name is the one used in stageForNum with the num
// argument. The associated target stage name is the one used in the stageForNum
// with the num+1 argument.
func mockLink(num int, sourceField string, targetField string) *link.Link {
	name := testutil.LinkNameForNum(num)
	sourceStage := testutil.LinkSourceStageForNum(num)
	targetStage := testutil.LinkTargetStageForNum(num)
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
		assert.NilError(t, s.flowManager.RegisterStage(st), "register stage")
	}
}

// populateLinks creates the links in the server, asserting any occurred errors.
func populateLinks(t *testing.T, s *Server, links []*link.Link) {
	// Bypass CreateLink verifications
	store := s.linkStore
	for _, l := range links {
		source, ok := s.stageStore.GetByName(l.SourceStage())
		assert.Assert(t, ok, "source does not exist")
		target, ok := s.stageStore.GetByName(l.TargetStage())
		assert.Assert(t, ok, "target does not exist")
		assert.NilError(t, store.Create(l), "populate with links")
		err := s.flowManager.RegisterLink(source, target, l)
		assert.NilError(t, err, "register link")
	}
}
