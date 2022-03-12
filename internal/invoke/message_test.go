package invoke

import (
	"github.com/DuarteMRAlves/maestro/internal"
	"github.com/DuarteMRAlves/maestro/test/protobuf/unit"
	"github.com/jhump/protoreflect/dynamic"
	"gotest.tools/v3/assert"
	"testing"
)

func TestNewFieldSetter(t *testing.T) {
	field := internal.NewMessageField("val")

	msg, err := NewDynamicMessage(&unit.DynamicTestMessage{})
	assert.NilError(t, err, "create dynamic message")

	err = msg.SetField(field, int32(1))
	assert.NilError(t, err, "set error")

	grpcMsg, err := dynamic.AsDynamicMessage(msg.GrpcMessage())
	assert.NilError(t, err, "dynamic grpc message")

	val, ok := grpcMsg.GetFieldByName(field.Unwrap()).(int32)
	assert.Assert(t, ok, "cast to int on grpc message")
	assert.Equal(t, int32(1), val)
}

func TestNewFieldGetter(t *testing.T) {
	field := internal.NewMessageField("inner")

	pbMsg := &unit.DynamicTestMessage{
		Inner: &unit.DynamicTestMessageInner{Val: "val"},
	}

	msg, err := NewDynamicMessage(pbMsg)
	assert.NilError(t, err, "create dynamic message")

	inner, err := msg.GetField(field)
	assert.NilError(t, err, "get inner error")

	innerDyn, err := dynamic.AsDynamicMessage(inner.GrpcMessage())
	assert.NilError(t, err, "inner as dynamic message")

	pbInner := &unit.DynamicTestMessageInner{}
	err = innerDyn.ConvertTo(pbInner)
	assert.NilError(t, err, "convert dynamic to inner")

	assert.Equal(t, "val", pbInner.Val)
}
