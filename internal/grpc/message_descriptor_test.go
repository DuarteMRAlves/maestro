package grpc

import (
	"github.com/DuarteMRAlves/maestro/internal"
	"github.com/DuarteMRAlves/maestro/test/protobuf/unit"
	"github.com/golang/protobuf/proto"
	"gotest.tools/v3/assert"
	"testing"
)

func TestMessageDescriptor_GetField(t *testing.T) {
	msgDesc, err := newMessageDescriptor(&unit.TestMessage1{})
	assert.NilError(t, err, "create message descriptor")

	field := internal.NewMessageField("field4")

	fieldDesc, err := msgDesc.GetField(field)
	assert.NilError(t, err, "get field")

	grpcFieldDesc, ok := fieldDesc.(messageDescriptor)
	assert.Assert(t, ok, "cast to message descriptor")
	assert.Equal(t, "InternalMessage1", grpcFieldDesc.desc.GetName())
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
	for name, test := range tests {
		t.Run(
			name,
			func(t *testing.T) {
				desc1, err := newMessageDescriptor(test.msg1)
				assert.NilError(t, err, "create message descriptor 1")
				desc2, err := newMessageDescriptor(test.msg2)
				assert.NilError(t, err, "create message descriptor 2")

				assert.Equal(t, test.expected, desc1.Compatible(desc2))
				assert.Equal(t, test.expected, desc2.Compatible(desc1))
			},
		)
	}
}
