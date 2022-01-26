package execution

import (
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/reflection"
	"github.com/DuarteMRAlves/maestro/internal/storage"
	"github.com/DuarteMRAlves/maestro/tests/pb"
	"github.com/jhump/protoreflect/desc"
	"gotest.tools/v3/assert"
	"reflect"
	"sync"
	"testing"
)

func TestManager_Register_NoFields(t *testing.T) {
	rpcManager := &reflection.MockManager{Rpcs: sync.Map{}}
	s1 := stage1(t, rpcManager)
	s2 := stage2(t, rpcManager)
	l := storage.NewLink("link-name", s1.Name(), "", s2.Name(), "")
	m := NewManager(rpcManager)
	err := m.RegisterLink(s1, s2, l)
	assert.NilError(t, err, "register error")
}

func TestManager_Register_WithFields(t *testing.T) {
	rpcManager := &reflection.MockManager{Rpcs: sync.Map{}}
	s1 := stage1(t, rpcManager)
	s2 := stage2(t, rpcManager)
	l := storage.NewLink(
		"link-name",
		s1.Name(),
		"field4",
		s2.Name(),
		"fieldName4",
	)
	m := NewManager(rpcManager)
	err := m.RegisterLink(s1, s2, l)
	assert.NilError(t, err, "register error")
}

func stage1(t *testing.T, rpcManager *reflection.MockManager) *storage.Stage {
	testMsg1Type := reflect.TypeOf(pb.TestMessage1{})

	testMsg1Desc, err := desc.LoadMessageDescriptorForType(testMsg1Type)
	assert.NilError(t, err, "load desc test message 1")

	message1, err := reflection.NewMessage(testMsg1Desc)
	assert.NilError(t, err, "test message 1")

	rpcManager.Rpcs.Store(
		api.StageName("stage-1"),
		&reflection.MockRPC{
			Name_: "rpc-1",
			FQN:   "service-1/rpc-1",
			In:    message1,
			Out:   message1,
		},
	)

	return storage.NewStage(
		"stage-1",
		storage.NewRpcSpec("address-1", "service-1", "rpc-1"),
		"asset-1",
		nil,
	)
}

func stage2(t *testing.T, rpcManager *reflection.MockManager) *storage.Stage {
	testMsg2Type := reflect.TypeOf(pb.TestMessageDiffNames{})

	testMsg2Desc, err := desc.LoadMessageDescriptorForType(testMsg2Type)
	assert.NilError(t, err, "load desc test message 2")

	message2, err := reflection.NewMessage(testMsg2Desc)
	assert.NilError(t, err, "test message 2")

	rpcManager.Rpcs.Store(
		api.StageName("stage-2"),
		&reflection.MockRPC{
			Name_: "rpc-2",
			FQN:   "service-2/rpc-2",
			In:    message2,
			Out:   message2,
		},
	)

	return storage.NewStage(
		"stage-2",
		storage.NewRpcSpec("address-2", "service-2", "rpc-2"),
		"asset-2",
		nil,
	)
}
