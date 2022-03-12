package errdefs

import (
	"fmt"
	"gotest.tools/v3/assert"
	"testing"
)

func TestAlreadyExistsWithError(t *testing.T) {
	var ok bool
	err := AlreadyExistsWithError(dummyErr)
	_, ok = err.(AlreadyExists)
	assert.Assert(t, ok, "AlreadyExists interface")
	_, ok = err.(alreadyExists)
	assert.Assert(t, ok, "alreadyExists struct")
	msg := err.Error()
	assert.Equal(t, dummyErrMsg, msg, "error message")
}

func TestAlreadyExistsWithMsg(t *testing.T) {
	var ok bool
	err := AlreadyExistsWithMsg(dummyErrMsg)
	_, ok = err.(AlreadyExists)
	assert.Assert(t, ok, "AlreadyExists interface")
	_, ok = err.(alreadyExists)
	assert.Assert(t, ok, "alreadyExists struct")
	msg := err.Error()
	assert.Equal(t, dummyErrMsg, msg, "error message")
}

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

func TestFailedPreconditionWithError(t *testing.T) {
	var ok bool
	err := FailedPreconditionWithError(dummyErr)
	_, ok = err.(FailedPrecondition)
	assert.Assert(t, ok, "FailedPrecondition interface")
	_, ok = err.(failedPrecondition)
	assert.Assert(t, ok, "failedPrecondition struct")
	msg := err.Error()
	assert.Equal(t, dummyErrMsg, msg, "error message")
}

func TestFailedPreconditionWithMsg(t *testing.T) {
	var ok bool
	err := FailedPreconditionWithMsg(dummyErrMsg)
	_, ok = err.(FailedPrecondition)
	assert.Assert(t, ok, "FailedPrecondition interface")
	_, ok = err.(failedPrecondition)
	assert.Assert(t, ok, "failedPrecondition struct")
	msg := err.Error()
	assert.Equal(t, dummyErrMsg, msg, "error message")
}

func TestUnavailableWithError(t *testing.T) {
	var ok bool
	err := UnavailableWithError(dummyErr)
	_, ok = err.(Unavailable)
	assert.Assert(t, ok, "Unavailable interface")
	_, ok = err.(unavailable)
	assert.Assert(t, ok, "unavailable struct")
	msg := err.Error()
	assert.Equal(t, dummyErrMsg, msg, "error message")
}

func TestUnavailableWithMsg(t *testing.T) {
	var ok bool
	err := UnavailableWithMsg(dummyErrMsg)
	_, ok = err.(Unavailable)
	assert.Assert(t, ok, "Unavailable interface")
	_, ok = err.(unavailable)
	assert.Assert(t, ok, "unavailable struct")
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
			name:      "already exists error",
			err:       AlreadyExistsWithMsg(dummyErrMsg),
			valTypeFn: IsAlreadyExists,
		},
		{
			name:      "invalid argument error",
			err:       InvalidArgumentWithMsg(dummyErrMsg),
			valTypeFn: IsInvalidArgument,
		},
		{
			name:      "failed precondition error",
			err:       FailedPreconditionWithMsg(dummyErrMsg),
			valTypeFn: IsFailedPrecondition,
		},
		{
			name:      "unavailable error",
			err:       UnavailableWithMsg(dummyErrMsg),
			valTypeFn: IsUnavailable,
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
