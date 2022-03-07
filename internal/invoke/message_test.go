package invoke

import (
	"github.com/DuarteMRAlves/maestro/internal/domain"
	"github.com/DuarteMRAlves/maestro/test/protobuf/unit"
	"github.com/jhump/protoreflect/dynamic"
	"gotest.tools/v3/assert"
	"testing"
)

func TestNewFieldSetter(t *testing.T) {
	field, err := domain.NewMessageField("val")
	assert.NilError(t, err, "create message field")

	setter := NewFieldSetter(field)

	emptyDyn, err := dynamic.AsDynamicMessage(&unit.DynamicTestMessage{})
	assert.NilError(t, err, "create dynamic message")

	empty := newDynamicMessage(emptyDyn)

	res := setter(empty, int32(1))
	assert.Assert(t, !res.IsError(), "error result")

	msg := res.Unwrap()
	grpcMsg, err := dynamic.AsDynamicMessage(msg.GrpcMessage())
	assert.NilError(t, err, "dynamic grpc message")

	val, ok := grpcMsg.GetFieldByName(field.Unwrap()).(int32)
	assert.Assert(t, ok, "cast to int on grpc message")
	assert.Equal(t, int32(1), val)

	emptyVal, ok := emptyDyn.GetFieldByName(field.Unwrap()).(int32)
	assert.Assert(t, ok, "cast to int on empty dynamic field")
	assert.Equal(t, int32(0), emptyVal)
}

func TestNewFieldGetter(t *testing.T) {
	field, err := domain.NewMessageField("inner")
	assert.NilError(t, err, "create message field")

	getter := NewFieldGetter(field)

	pbMsg := &unit.DynamicTestMessage{
		Inner: &unit.DynamicTestMessageInner{Val: "val"},
	}

	dyn, err := dynamic.AsDynamicMessage(pbMsg)
	assert.NilError(t, err, "dynamic message from proto")

	msg := newDynamicMessage(dyn)

	res := getter(msg)
	assert.Assert(t, !res.IsError(), "get result error")

	inner := res.Unwrap()
	innerDyn, err := dynamic.AsDynamicMessage(inner.GrpcMessage())
	assert.NilError(t, err, "inner as dynamic message")

	pbInner := &unit.DynamicTestMessageInner{}
	err = innerDyn.ConvertTo(pbInner)
	assert.NilError(t, err, "convert dynamic to inner")

	assert.Equal(t, "val", pbInner.Val)
}
