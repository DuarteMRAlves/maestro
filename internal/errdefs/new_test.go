package errdefs

import (
	"fmt"
	"gotest.tools/v3/assert"
	"testing"
)

func TestInvalidArgumentWithError(t *testing.T) {
	var ok bool
	err := InvalidArgumentWithError(dummyErr)
	_, ok = err.(InvalidArgument)
	assert.Assert(t, ok, "InvalidArgument interface")
	_, ok = err.(invalidArgument)
	assert.Assert(t, ok, "invalidArgument struct")
	msg := err.Error()
	assert.Equal(t, dummyErrMsg, msg, "error message")
}

func TestInvalidArgumentWithMessage(t *testing.T) {
	var ok bool
	err := InvalidArgumentWithMsg(dummyErrMsg)
	_, ok = err.(InvalidArgument)
	assert.Assert(t, ok, "InvalidArgument interface")
	_, ok = err.(invalidArgument)
	assert.Assert(t, ok, "invalidArgument struct")
	msg := err.Error()
	assert.Equal(t, dummyErrMsg, msg, "error message")
}

func TestInternalWithError(t *testing.T) {
	var ok bool
	err := InternalWithError(dummyErr)
	_, ok = err.(Internal)
	assert.Assert(t, ok, "Internal interface")
	_, ok = err.(internal)
	assert.Assert(t, ok, "internal struct")
	msg := err.Error()
	assert.Equal(t, dummyErrMsg, msg, "error message")
}

func TestInternalWithMessage(t *testing.T) {
	var ok bool
	err := InternalWithMsg(dummyErrMsg)
	_, ok = err.(Internal)
	assert.Assert(t, ok, "Internal interface")
	_, ok = err.(internal)
	assert.Assert(t, ok, "internal struct")
	msg := err.Error()
	assert.Equal(t, dummyErrMsg, msg, "error message")
}

func TestUnknownWithError(t *testing.T) {
	var ok bool
	err := UnknownWithError(dummyErr)
	_, ok = err.(Unknown)
	assert.Assert(t, ok, "Unknown interface")
	_, ok = err.(unknown)
	assert.Assert(t, ok, "unknown struct")
	msg := err.Error()
	assert.Equal(t, dummyErrMsg, msg, "error message")
}

func TestUnknownWithMessage(t *testing.T) {
	var ok bool
	err := UnknownWithMsg(dummyErrMsg)
	_, ok = err.(Unknown)
	assert.Assert(t, ok, "Unknown interface")
	_, ok = err.(unknown)
	assert.Assert(t, ok, "unknown struct")
	msg := err.Error()
	assert.Equal(t, dummyErrMsg, msg, "error message")
}

func TestPrependMsg(t *testing.T) {
	prependFormat := "prepend message with %s, %s"
	prependArgs := []interface{}{"arg1", "arg2"}

	prependMsg := fmt.Sprintf(prependFormat, prependArgs...)
	expectedMsg := fmt.Sprintf("%s: %s", prependMsg, dummyErrMsg)

	tests := []struct {
		name      string
		err       error
		valTypeFn func(err error) bool
	}{
		{
			name:      "invalid argument error",
			err:       InvalidArgumentWithMsg(dummyErrMsg),
			valTypeFn: IsInvalidArgument,
		},
		{
			name:      "internal error",
			err:       InternalWithMsg(dummyErrMsg),
			valTypeFn: IsInternal,
		},
		{
			name:      "unknown error",
			err:       UnknownWithMsg(dummyErrMsg),
			valTypeFn: IsUnknown,
		},
		{
			name:      "fmt error",
			err:       fmt.Errorf(dummyErrMsg),
			valTypeFn: IsUnknown,
		},
	}

	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				wrappedErr := PrependMsg(
					test.err,
					prependFormat,
					prependArgs...,
				)
				assert.Assert(t, test.valTypeFn(wrappedErr), "correct type")
				msg := wrappedErr.Error()
				assert.Equal(t, expectedMsg, msg, "correct message")
			},
		)
	}
}
