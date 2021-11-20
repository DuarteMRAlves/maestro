package errdefs

import (
	"github.com/DuarteMRAlves/maestro/internal/assert"
	"testing"
)

func TestAlreadyExistsWithError(t *testing.T) {
	var ok bool
	err := AlreadyExistsWithError(dummyErr)
	_, ok = err.(AlreadyExists)
	assert.IsTrue(t, ok, "AlreadyExists interface")
	_, ok = err.(alreadyExists)
	assert.IsTrue(t, ok, "alreadyExists struct")
	msg := err.Error()
	assert.DeepEqual(t, dummyErrMsg, msg, "error message")
}

func TestAlreadyExistsWithMsg(t *testing.T) {
	var ok bool
	err := AlreadyExistsWithMsg(dummyErrMsg)
	_, ok = err.(AlreadyExists)
	assert.IsTrue(t, ok, "AlreadyExists interface")
	_, ok = err.(alreadyExists)
	assert.IsTrue(t, ok, "alreadyExists struct")
	msg := err.Error()
	assert.DeepEqual(t, dummyErrMsg, msg, "error message")
}

func TestNotFoundWithError(t *testing.T) {
	var ok bool
	err := NotFoundWithError(dummyErr)
	_, ok = err.(NotFound)
	assert.IsTrue(t, ok, "NotFound interface")
	_, ok = err.(notFound)
	assert.IsTrue(t, ok, "notFound struct")
	msg := err.Error()
	assert.DeepEqual(t, dummyErrMsg, msg, "error message")
}

func TestNotFoundWithMessage(t *testing.T) {
	var ok bool
	err := NotFoundWithMsg(dummyErrMsg)
	_, ok = err.(NotFound)
	assert.IsTrue(t, ok, "NotFound interface")
	_, ok = err.(notFound)
	assert.IsTrue(t, ok, "notFound struct")
	msg := err.Error()
	assert.DeepEqual(t, dummyErrMsg, msg, "error message")
}

func TestUnknownWithError(t *testing.T) {
	var ok bool
	err := UnknownWithError(dummyErr)
	_, ok = err.(Unknown)
	assert.IsTrue(t, ok, "Unknown interface")
	_, ok = err.(unknown)
	assert.IsTrue(t, ok, "unknown struct")
	msg := err.Error()
	assert.DeepEqual(t, dummyErrMsg, msg, "error message")
}

func TestUnknownWithMessage(t *testing.T) {
	var ok bool
	err := UnknownWithMsg(dummyErrMsg)
	_, ok = err.(Unknown)
	assert.IsTrue(t, ok, "Unknown interface")
	_, ok = err.(unknown)
	assert.IsTrue(t, ok, "unknown struct")
	msg := err.Error()
	assert.DeepEqual(t, dummyErrMsg, msg, "error message")
}
