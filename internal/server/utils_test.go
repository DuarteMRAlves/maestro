package server

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/asset"
	"github.com/DuarteMRAlves/maestro/internal/orchestration"
	"github.com/DuarteMRAlves/maestro/internal/reflection"
	"github.com/DuarteMRAlves/maestro/internal/testutil"
	mockreflection "github.com/DuarteMRAlves/maestro/internal/testutil/mock/reflection"
	"github.com/dgraph-io/badger/v3"
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
	rpcManager *mockreflection.Manager,
) *orchestration.Stage {
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

	rpcManager.Rpcs.Store(
		testutil.StageNameForNum(num),
		&mockreflection.RPC{
			Name_: testutil.StageRpcForNum(num),
			FQN: fmt.Sprintf(
				"%s/%s",
				testutil.StageServiceForNum(num),
				testutil.StageRpcForNum(num),
			),
			In:    reqMsg,
			Out:   resMsg,
			Unary: true,
		},
	)

	return orchestration.NewStage(
		testutil.StageNameForNum(num),
		orchestration.NewRpcSpec(
			testutil.StageAddressForNum(num),
			testutil.StageServiceForNum(num),
			testutil.StageRpcForNum(num),
		),
		testutil.AssetNameForNum(num),
		nil,
	)
}

// populateStages creates the stages in the server, asserting any occurred
// errors.
func populateStages(
	t *testing.T,
	s *Server,
	txn *badger.Txn,
	stages []*orchestration.Stage,
) {
	for _, st := range stages {
		orchestration.PersistStage(txn, st)
		assert.NilError(t, s.flowManager.RegisterStage(st), "register stage")
	}
}
