package grpc

import (
	"testing"

	"github.com/DuarteMRAlves/maestro/internal/message"
	"github.com/DuarteMRAlves/maestro/test/protobuf/unit"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jhump/protoreflect/dynamic"
)

func TestNewFieldSetter(t *testing.T) {
	field := message.Field("inner")
	pbInner := &unit.TestMessageInner{Val: "val"}

	msg, err := newMessage(&unit.TestMessage{})
	if err != nil {
		t.Fatalf("create outer message: %s", err)
	}

	inner, err := newMessage(pbInner)
	if err != nil {
		t.Fatalf("create inner message: %s", err)
	}

	err = msg.Set(field, inner)
	if err != nil {
		t.Fatalf("set field error: %s", err)
	}

	fieldDyn, ok := msg.dynMsg.GetFieldByName(string(field)).(*dynamic.Message)
	if !ok {
		t.Fatalf("unable to cast to dynamic message on fieldDyn")
	}

	pbField := &unit.TestMessageInner{}
	err = fieldDyn.ConvertTo(pbField)
	if err != nil {
		t.Fatalf("convert fieldDyn to pbField: %s", err)
	}

	cmpOpts := cmpopts.IgnoreUnexported(unit.TestMessageInner{})
	if diff := cmp.Diff(pbInner, pbField, cmpOpts); diff != "" {
		t.Fatalf("mismatch on original and retreived:\n%s", diff)
	}
}

func TestNewFieldGetter(t *testing.T) {
	field := message.Field("inner")

	pbInner := &unit.TestMessageInner{Val: "val"}
	pbMsg := &unit.TestMessage{Inner: pbInner}

	msg, err := newMessage(pbMsg)
	if err != nil {
		t.Fatalf("create message: %s", err)
	}

	fieldVal, err := msg.Get(field)
	if err != nil {
		t.Fatalf("get field: %s", err)
	}

	fieldGrpc, ok := fieldVal.(*messageInstance)
	if !ok {
		t.Fatalf("unable to cast to fieldVal to fieldGrpc")
	}

	pbField := &unit.TestMessageInner{}
	err = fieldGrpc.dynMsg.ConvertTo(pbField)
	if err != nil {
		t.Fatalf("convert fieldGrpc to pbField: %s", err)
	}

	cmpOpts := cmpopts.IgnoreUnexported(unit.TestMessageInner{})
	if diff := cmp.Diff(pbInner, pbField, cmpOpts); diff != "" {
		t.Fatalf("mismatch on original and retreived:\n%s", diff)
	}
}
