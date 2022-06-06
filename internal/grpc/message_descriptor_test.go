package grpc

import (
	"testing"

	"github.com/DuarteMRAlves/maestro/internal/compiled"
	"github.com/DuarteMRAlves/maestro/test/protobuf/unit"
	"github.com/golang/protobuf/proto"
	"github.com/google/go-cmp/cmp"
)

func TestMessageDescriptor_GetField(t *testing.T) {
	msgDesc, err := newMessageDescriptor(&unit.TestMessage1{})
	if err != nil {
		t.Fatalf("create message descriptor: %s", err)
	}

	field := compiled.NewMessageField("field4")

	fieldDesc, err := msgDesc.GetField(field)
	if err != nil {
		t.Fatalf("get field: %s", err)
	}

	grpcFieldDesc, ok := fieldDesc.(messageDescriptor)
	if !ok {
		t.Fatalf("unable to cast to fieldDesc to grpcFieldDesc")
	}
	if diff := cmp.Diff("InternalMessage1", grpcFieldDesc.desc.GetName()); diff != "" {
		t.Fatalf("mismatch on original and retreived names:\n%s", diff)
	}
}

func TestCompatibleDescriptors(t *testing.T) {
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
				desc1, err := newMessageDescriptor(tc.msg1)
				if err != nil {
					t.Fatalf("create message descriptor 1: %s", err)
				}
				desc2, err := newMessageDescriptor(tc.msg2)
				if err != nil {
					t.Fatalf("create message descriptor 2: %s", err)
				}

				if diff := cmp.Diff(tc.expected, desc1.Compatible(desc2)); diff != "" {
					t.Fatalf("mismatch on desc1.Compatible(desc2):\n%s", diff)
				}
				if diff := cmp.Diff(tc.expected, desc2.Compatible(desc1)); diff != "" {
					t.Fatalf("mismatch on desc2.Compatible(desc1):\n%s", diff)
				}
			},
		)
	}
}
