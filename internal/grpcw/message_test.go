package grpcw

import (
	"reflect"
	"testing"

	"github.com/DuarteMRAlves/maestro/internal/message"
	"github.com/DuarteMRAlves/maestro/test/protobuf/unit"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestInstanceSet(t *testing.T) {
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

func TestInstanceGet(t *testing.T) {
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

func TestTypeSubfield(t *testing.T) {
	pbOuter := &unit.TestMessage1{}
	typeOuter := messageType{pbOuter.ProtoReflect().Type()}
	field := message.Field("field4")

	subfield, err := typeOuter.Subfield(field)
	if err != nil {
		t.Fatalf("subfield: %s", err)
	}
	innerType, ok := subfield.(messageType)
	if !ok {
		t.Fatalf("unable to cast to subfield to messageType")
	}
	exp := protoreflect.Name("InternalMessage1")
	actual := innerType.t.Descriptor().Name()
	if diff := cmp.Diff(exp, actual); diff != "" {
		t.Fatalf("names mismatch:\n%s", diff)
	}
}

func TestTypeCompatible(t *testing.T) {
	tests := map[string]struct {
		msg1     proto.Message
		msg2     proto.Message
		expected bool
	}{
		"equal messages": {
			msg1:     &unit.TestMessage1{},
			msg2:     &unit.TestMessage1{},
			expected: true,
		},
		"different message names": {
			msg1:     &unit.TestMessage1{},
			msg2:     &unit.TestMessage2{},
			expected: true,
		},
		"different field names": {
			msg1:     &unit.TestMessage1{},
			msg2:     &unit.TestMessageDiffNames{},
			expected: true,
		},
		"different non common fields": {
			msg1:     &unit.TestMessage1{},
			msg2:     &unit.TestMessageDiffFields{},
			expected: true,
		},
		"outer message wrong cardinality": {
			msg1:     &unit.TestMessage1{},
			msg2:     &unit.TestWrongOuterCardinality{},
			expected: false,
		},
		"inner message wrong cardinality": {
			msg1:     &unit.TestMessage1{},
			msg2:     &unit.TestWrongInnerCardinality{},
			expected: false,
		},
		"outer message wrong field type": {
			msg1:     &unit.TestMessage1{},
			msg2:     &unit.TestWrongOuterFieldType{},
			expected: false,
		},
		"inner message wrong field type": {
			msg1:     &unit.TestMessage1{},
			msg2:     &unit.TestWrongInnerFieldType{},
			expected: false,
		},
	}
	for name, tc := range tests {
		t.Run(
			name,
			func(t *testing.T) {
				t1 := messageType{tc.msg1.ProtoReflect().Type()}
				t2 := messageType{tc.msg2.ProtoReflect().Type()}

				if diff := cmp.Diff(tc.expected, t1.Compatible(t2)); diff != "" {
					t.Fatalf("mismatch on t1.Compatible(t2):\n%s", diff)
				}
				if diff := cmp.Diff(tc.expected, t2.Compatible(t1)); diff != "" {
					t.Fatalf("mismatch on t2.Compatible(t1):\n%s", diff)
				}
			},
		)
	}
}
