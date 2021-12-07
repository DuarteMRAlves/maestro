package errdefs

import (
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

func TestNotFoundWithError(t *testing.T) {
	var ok bool
	err := NotFoundWithError(dummyErr)
	_, ok = err.(NotFound)
	assert.Assert(t, ok, "NotFound interface")
	_, ok = err.(notFound)
	assert.Assert(t, ok, "notFound struct")
	msg := err.Error()
	assert.Equal(t, dummyErrMsg, msg, "error message")
}

func TestNotFoundWithMessage(t *testing.T) {
	var ok bool
	err := NotFoundWithMsg(dummyErrMsg)
	_, ok = err.(NotFound)
	assert.Assert(t, ok, "NotFound interface")
	_, ok = err.(notFound)
	assert.Assert(t, ok, "notFound struct")
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
