package server

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/reflection"
	"github.com/DuarteMRAlves/maestro/internal/storage"
	"github.com/DuarteMRAlves/maestro/internal/testutil"
	mockreflection "github.com/DuarteMRAlves/maestro/internal/testutil/mock/reflection"
	"github.com/dgraph-io/badger/v3"
	"github.com/jhump/protoreflect/desc"
	"gotest.tools/v3/assert"
	"reflect"
	"testing"
)

// assetForNum deterministically creates an asset with the given number.
func assetForNum(num int) *storage.Asset {
	name := testutil.AssetNameForNum(num)
	img := testutil.AssetImageForNum(num)
	return storage.NewAsset(name, img)
}

// mockStage deterministically creates a stage with the given number.
// The associated asset name is the one used in assetForNum.
func mockStage(
	t *testing.T,
	num int,
	req interface{},
	res interface{},
	rpcManager *mockreflection.Manager,
) *storage.Stage {
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

	return storage.NewStage(
		testutil.StageNameForNum(num),
		storage.NewRpcSpec(
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
	stages []*storage.Stage,
) {
	for _, st := range stages {
		storage.PersistStage(txn, st)
		assert.NilError(t, s.flowManager.RegisterStage(st), "register stage")
	}
}
