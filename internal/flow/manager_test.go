package flow

import (
	"github.com/DuarteMRAlves/maestro/internal/orchestration"
	"github.com/DuarteMRAlves/maestro/internal/reflection"
	mockreflection "github.com/DuarteMRAlves/maestro/internal/testutil/mock/reflection"
	"github.com/DuarteMRAlves/maestro/tests/pb"
	"github.com/jhump/protoreflect/desc"
	"gotest.tools/v3/assert"
	"reflect"
	"testing"
)

func TestManager_Register_NoFields(t *testing.T) {
	s1 := stage1(t)
	s2 := stage2(t)
	l := orchestration.NewLink("link-name", s1.Name(), "", s2.Name(), "")
	manager := NewManager()
	err := manager.RegisterLink(s1, s2, l)
	assert.NilError(t, err, "register error")
}

func TestManager_Register_WithFields(t *testing.T) {
	s1 := stage1(t)
	s2 := stage2(t)
	l := orchestration.NewLink("link-name", s1.Name(), "field4", s2.Name(), "fieldName4")
	manager := NewManager()
	err := manager.RegisterLink(s1, s2, l)
	assert.NilError(t, err, "register error")
}

func stage1(t *testing.T) *orchestration.Stage {
	testMsg1Type := reflect.TypeOf(pb.TestMessage1{})

	testMsg1Desc, err := desc.LoadMessageDescriptorForType(testMsg1Type)
	assert.NilError(t, err, "load desc test message 1")

	message1, err := reflection.NewMessage(testMsg1Desc)
	assert.NilError(t, err, "test message 1")

	return orchestration.NewStage(
		"stage-1",
		"asset-1",
		"address-1",
		&mockreflection.RPC{
			Name_: "rpc-1",
			FQN:   "service-1/rpc-1",
			In:    message1,
			Out:   message1,
		},
		nil)
}

func stage2(t *testing.T) *orchestration.Stage {
	testMsg2Type := reflect.TypeOf(pb.TestMessageDiffNames{})

	testMsg2Desc, err := desc.LoadMessageDescriptorForType(testMsg2Type)
	assert.NilError(t, err, "load desc test message 2")

	message2, err := reflection.NewMessage(testMsg2Desc)
	assert.NilError(t, err, "test message 2")

	return orchestration.NewStage(
		"stage-2",
		"asset-2",
		"address-2",
		&mockreflection.RPC{
			Name_: "rpc-2",
			FQN:   "service-2/rpc-2",
			In:    message2,
			Out:   message2,
		},
		nil)
}

func incompatibleStage(t *testing.T) *orchestration.Stage {
	testIncompatibleType := reflect.TypeOf(pb.TestWrongOuterFieldType{})

	testIncompatibleDesc, err := desc.LoadMessageDescriptorForType(
		testIncompatibleType)
	assert.NilError(t, err, "load desc test message incompatible")

	messageIncompatible, err := reflection.NewMessage(testIncompatibleDesc)
	assert.NilError(t, err, "test message incompatible")

	return orchestration.NewStage(
		"stage-3",
		"asset-incompatible",
		"address-incompatible",
		&mockreflection.RPC{
			Name_: "rpc-incompatible",
			FQN:   "service-2/rpc-incompatible",
			In:    messageIncompatible,
			Out:   messageIncompatible,
		},
		nil)
}
