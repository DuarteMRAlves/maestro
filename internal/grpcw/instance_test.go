package grpcw

import (
	"reflect"
	"testing"

	"github.com/DuarteMRAlves/maestro/internal/message"
	"github.com/DuarteMRAlves/maestro/test/protobuf/unit"
	"github.com/google/go-cmp/cmp"
)

func TestSet(t *testing.T) {
	field := message.Field("inner")
	val := "val"

	pbOuter := &unit.TestMessage{}
	pbInner := &unit.TestMessageInner{Val: val}

	msgOuter := messageInstance{m: pbOuter.ProtoReflect()}
	msgInner := messageInstance{m: pbInner.ProtoReflect()}

	err := msgOuter.Set(field, msgInner)
	if err != nil {
		t.Fatalf("set field error: %s", err)
	}

	iface := msgOuter.m.Interface()
	unwrapped, ok := iface.(*unit.TestMessage)
	if !ok {
		t.Fatalf("unwrapped type mismatch: expected *unit.TestMessage, got %s", reflect.TypeOf(iface))
	}
	if diff := cmp.Diff(val, unwrapped.Inner.Val); diff != "" {
		t.Fatalf("unwrapped val mismatch:\n%s", diff)
	}
}

func TestGet(t *testing.T) {
	field := message.Field("inner")
	val := "val"

	pbInner := &unit.TestMessageInner{Val: val}
	pbOuter := &unit.TestMessage{Inner: pbInner}

	msgOuter := messageInstance{m: pbOuter.ProtoReflect()}

	fieldVal, err := msgOuter.Get(field)
	if err != nil {
		t.Fatalf("get field: %s", err)
	}

	msgInner, ok := fieldVal.(messageInstance)
	if !ok {
		t.Fatalf("fieldVal type mismatch: expected messageInstance, got %s", reflect.TypeOf(fieldVal))
	}

	iface := msgInner.m.Interface()
	unwrapped, ok := iface.(*unit.TestMessageInner)
	if !ok {
		t.Fatalf("unwrapped type mismatch: expected *unit.TestMessage, got %s", reflect.TypeOf(iface))
	}
	if diff := cmp.Diff(val, unwrapped.Val); diff != "" {
		t.Fatalf("unwrapped val mismatch:\n%s", diff)
	}
}
