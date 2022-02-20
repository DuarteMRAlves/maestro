package rpc

import (
	"github.com/jhump/protoreflect/desc"
	"gotest.tools/v3/assert"
	"testing"
)

func TestMessage_Compatible_True(t *testing.T) {
	tests := []struct {
		name     string
		message1 string
		message2 string
	}{
		{
			name:     "equal names",
			message1: "unit.TestMessage1",
			message2: "unit.TestMessage2",
		},
		{
			name:     "equal descriptors different message names",
			message1: "unit.TestMessage1",
			message2: "unit.TestMessage2",
		},
		{
			name:     "equal descriptors different field names",
			message1: "unit.TestMessage1",
			message2: "unit.TestMessageDiffNames",
		},
		{
			name:     "equal descriptors different fields",
			message1: "unit.TestMessage1",
			message2: "unit.TestMessageDiffFields",
		},
	}
	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				descriptor1, err := desc.LoadMessageDescriptor(test.message1)
				assert.NilError(t, err, "load message descriptor 1")

				descriptor2, err := desc.LoadMessageDescriptor(test.message2)
				assert.NilError(t, err, "load message descriptor 2")

				message1 := newMessageInternal(descriptor1)
				message2 := newMessageInternal(descriptor2)

				assert.Assert(
					t,
					message1.Compatible(message2),
					"messages compatible is not true",
				)
			},
		)
	}
}

func TestMessage_Compatible_False(t *testing.T) {
	tests := []struct {
		name     string
		message1 string
		message2 string
	}{
		{
			name:     "outer message wrong cardinality",
			message1: "unit.TestMessage1",
			message2: "unit.TestWrongOuterCardinality",
		},
		{
			name:     "inner message wrong cardinality",
			message1: "unit.TestMessage1",
			message2: "unit.TestWrongInnerCardinality",
		},
		{
			name:     "outer message wrong field type",
			message1: "unit.TestMessage1",
			message2: "unit.TestWrongOuterFieldType",
		},
		{
			name:     "inner message wrong field type",
			message1: "unit.TestMessage1",
			message2: "unit.TestWrongInnerFieldType",
		},
	}
	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				descriptor1, err := desc.LoadMessageDescriptor(test.message1)
				assert.NilError(t, err, "load message descriptor 1")
				assert.Assert(t, descriptor1 != nil)

				descriptor2, err := desc.LoadMessageDescriptor(test.message2)
				assert.NilError(t, err, "load message descriptor 2")
				assert.Assert(t, descriptor2 != nil)

				message1 := newMessageInternal(descriptor1)
				message2 := newMessageInternal(descriptor2)

				assert.Assert(
					t,
					!message1.Compatible(message2),
					"messages compatible is not false",
				)
			},
		)
	}
}
