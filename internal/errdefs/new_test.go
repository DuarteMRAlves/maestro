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
