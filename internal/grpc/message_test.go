package grpc

import (
	"github.com/DuarteMRAlves/maestro/internal"
	"github.com/DuarteMRAlves/maestro/test/protobuf/unit"
	"github.com/jhump/protoreflect/dynamic"
	"gotest.tools/v3/assert"
	"testing"
)

func TestNewFieldSetter(t *testing.T) {
	field := internal.NewMessageField("inner")

	msg, err := newMessage(&unit.TestMessage{})
	assert.NilError(t, err, "create outer message")

	inner, err := newMessage(&unit.TestMessageInner{Val: "val"})
	assert.NilError(t, err, "create inner message")

	err = msg.SetField(field, inner)
	assert.NilError(t, err, "set error")

	fieldVal, ok := msg.dynMsg.GetFieldByName(field.Unwrap()).(*dynamic.Message)
	assert.Assert(t, ok, "cast to dynamic message on inner message")

	pbInner := &unit.TestMessageInner{}
	err = fieldVal.ConvertTo(pbInner)
	assert.NilError(t, err, "convert dynamic to inner")
	assert.Equal(t, "val", pbInner.Val)
}

func TestNewFieldGetter(t *testing.T) {
	field := internal.NewMessageField("inner")

	pbMsg := &unit.TestMessage{Inner: &unit.TestMessageInner{Val: "val"}}

	msg, err := newMessage(pbMsg)
	assert.NilError(t, err, "create dynamic message")

	inner, err := msg.GetField(field)
	assert.NilError(t, err, "get inner error")

	innerGrpc, ok := inner.(*message)
	assert.Assert(t, ok, "cast inner to *grpc.message")

	pbInner := &unit.TestMessageInner{}
	err = innerGrpc.dynMsg.ConvertTo(pbInner)
	assert.NilError(t, err, "convert dynamic to inner")

	assert.Equal(t, "val", pbInner.Val)
}
