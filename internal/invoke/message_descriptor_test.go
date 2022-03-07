package invoke

import (
	"github.com/DuarteMRAlves/maestro/internal/domain"
	"github.com/DuarteMRAlves/maestro/test/protobuf/unit"
	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
	"gotest.tools/v3/assert"
	"reflect"
	"testing"
)

func TestMessageDescriptor_MessageFields(t *testing.T) {
	msgType := reflect.TypeOf(&unit.TestMessage1{})

	innerDesc, err := desc.LoadMessageDescriptorForType(msgType)
	assert.NilError(t, err, "load message descriptor")

	msgDesc, err := newMessageDescriptor(innerDesc)
	assert.NilError(t, err, "create message descriptor")

	fields := msgDesc.MessageFields()
	assert.Equal(t, 1, len(fields))

	field, err := domain.NewMessageField("field4")
	assert.NilError(t, err, "create inner field name")
	fieldDesc, ok := fields[field]
	assert.Assert(t, ok, "field4 exists")

	innerFieldDesc, ok := fieldDesc.(messageDescriptor)
	assert.Assert(t, ok, "cast for inner field desc")

	assert.Equal(t, "InternalMessage1", innerFieldDesc.desc.GetName())
}

func TestCompatibleDescriptors(t *testing.T) {
	tests := []struct {
		name     string
		msg1     proto.Message
		msg2     proto.Message
		expected bool
	}{
		{
			name:     "equal messages",
			msg1:     &unit.TestMessage1{},
			msg2:     &unit.TestMessage1{},
			expected: true,
		},
		{
			name:     "different message names",
			msg1:     &unit.TestMessage1{},
			msg2:     &unit.TestMessage2{},
			expected: true,
		},
		{
			name:     "different field names",
			msg1:     &unit.TestMessage1{},
			msg2:     &unit.TestMessageDiffNames{},
			expected: true,
		},
		{
			name:     "different non common fields",
			msg1:     &unit.TestMessage1{},
			msg2:     &unit.TestMessageDiffFields{},
			expected: true,
		},
		{
			name:     "outer message wrong cardinality",
			msg1:     &unit.TestMessage1{},
			msg2:     &unit.TestWrongOuterCardinality{},
			expected: false,
		},
		{
			name:     "inner message wrong cardinality",
			msg1:     &unit.TestMessage1{},
			msg2:     &unit.TestWrongInnerCardinality{},
			expected: false,
		},
		{
			name:     "outer message wrong field type",
			msg1:     &unit.TestMessage1{},
			msg2:     &unit.TestWrongOuterFieldType{},
			expected: false,
		},
		{
			name:     "inner message wrong field type",
			msg1:     &unit.TestMessage1{},
			msg2:     &unit.TestWrongInnerFieldType{},
			expected: false,
		},
	}
	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				type1 := reflect.TypeOf(test.msg1)
				innerDesc1, err := desc.LoadMessageDescriptorForType(type1)
				assert.NilError(t, err, "load message descriptor 1")

				type2 := reflect.TypeOf(test.msg2)
				innerDesc2, err := desc.LoadMessageDescriptorForType(type2)
				assert.NilError(t, err, "load message descriptor 2")

				desc1, err := newMessageDescriptor(innerDesc1)
				assert.NilError(t, err, "create message descriptor 1")
				desc2, err := newMessageDescriptor(innerDesc2)
				assert.NilError(t, err, "create message descriptor 2")

				assert.Equal(
					t,
					test.expected,
					CompatibleDescriptors(desc1, desc2),
				)
			},
		)
	}
}
