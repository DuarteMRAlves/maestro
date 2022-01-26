package server

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/rpc"
	"github.com/DuarteMRAlves/maestro/internal/storage"
	"github.com/DuarteMRAlves/maestro/internal/util"
	"github.com/dgraph-io/badger/v3"
	"github.com/jhump/protoreflect/desc"
	"gotest.tools/v3/assert"
	"reflect"
	"testing"
)

// assetForNum deterministically creates an asset with the given number.
func assetForNum(num int) *api.Asset {
	name := util.AssetNameForNum(num)
	img := util.AssetImageForNum(num)
	return &api.Asset{
		Name:  name,
		Image: img,
	}
}

// mockStage deterministically creates a stage with the given number.
// The associated asset name is the one used in assetForNum.
func mockStage(
	t *testing.T,
	num int,
	req interface{},
	res interface{},
	rpcManager *rpc.MockManager,
) *api.Stage {
	reqType := reflect.TypeOf(req)

	reqDesc, err := desc.LoadMessageDescriptorForType(reqType)
	assert.NilError(t, err, fmt.Sprintf("load req desc for stage: %d\n", num))

	reqMsg, err := rpc.NewMessage(reqDesc)
	assert.NilError(t, err, fmt.Sprintf("load req msg for stage: %d\n", num))

	resType := reflect.TypeOf(res)

	resDesc, err := desc.LoadMessageDescriptorForType(resType)
	assert.NilError(t, err, fmt.Sprintf("load res desc for stage: %d\n", num))

	resMsg, err := rpc.NewMessage(resDesc)
	assert.NilError(t, err, fmt.Sprintf("load res desc for stage: %d\n", num))

	rpcManager.Rpcs.Store(
		util.StageNameForNum(num),
		&rpc.MockRPC{
			Name_: util.StageRpcForNum(num),
			FQN: fmt.Sprintf(
				"%s/%s",
				util.StageServiceForNum(num),
				util.StageRpcForNum(num),
			),
			In:    reqMsg,
			Out:   resMsg,
			Unary: true,
		},
	)

	return &api.Stage{
		Name:    util.StageNameForNum(num),
		Phase:   api.StagePending,
		Service: util.StageServiceForNum(num),
		Rpc:     util.StageRpcForNum(num),
		Address: util.StageAddressForNum(num),
		Asset:   util.AssetNameForNum(num),
	}
}

// populateStages creates the stages in the server, asserting any occurred
// errors.
func populateStages(
	t *testing.T,
	s *Server,
	txn *badger.Txn,
	stages []*api.Stage,
) {
	helper := storage.NewTxnHelper(txn)
	for _, st := range stages {
		assert.NilError(t, helper.SaveStage(st), "persist stage")
		assert.NilError(t, s.flowManager.RegisterStage(st), "register stage")
	}
}
