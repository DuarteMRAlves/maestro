package errdefs

import (
	"github.com/DuarteMRAlves/maestro/internal/test"
	"testing"
)

func TestAlreadyExistsWithError(t *testing.T) {
	var ok bool
	err := AlreadyExistsWithError(dummyErr)
	_, ok = err.(AlreadyExists)
	test.IsTrue(t, ok, "AlreadyExists interface")
	_, ok = err.(alreadyExists)
	test.IsTrue(t, ok, "alreadyExists struct")
	msg := err.Error()
	test.DeepEqual(t, dummyErrMsg, msg, "error message")
}

func TestAlreadyExistsWithMsg(t *testing.T) {
	var ok bool
	err := AlreadyExistsWithMsg(dummyErrMsg)
	_, ok = err.(AlreadyExists)
	test.IsTrue(t, ok, "AlreadyExists interface")
	_, ok = err.(alreadyExists)
	test.IsTrue(t, ok, "alreadyExists struct")
	msg := err.Error()
	test.DeepEqual(t, dummyErrMsg, msg, "error message")
}

func TestNotFoundWithError(t *testing.T) {
	var ok bool
	err := NotFoundWithError(dummyErr)
	_, ok = err.(NotFound)
	test.IsTrue(t, ok, "NotFound interface")
	_, ok = err.(notFound)
	test.IsTrue(t, ok, "notFound struct")
	msg := err.Error()
	test.DeepEqual(t, dummyErrMsg, msg, "error message")
}

func TestNotFoundWithMessage(t *testing.T) {
	var ok bool
	err := NotFoundWithMsg(dummyErrMsg)
	_, ok = err.(NotFound)
	test.IsTrue(t, ok, "NotFound interface")
	_, ok = err.(notFound)
	test.IsTrue(t, ok, "notFound struct")
	msg := err.Error()
	test.DeepEqual(t, dummyErrMsg, msg, "error message")
}

func TestInvalidArgumentWithError(t *testing.T) {
	var ok bool
	err := InvalidArgumentWithError(dummyErr)
	_, ok = err.(InvalidArgument)
	test.IsTrue(t, ok, "InvalidArgument interface")
	_, ok = err.(invalidArgument)
	test.IsTrue(t, ok, "invalidArgument struct")
	msg := err.Error()
	test.DeepEqual(t, dummyErrMsg, msg, "error message")
}

func TestInvalidArgumentWithMessage(t *testing.T) {
	var ok bool
	err := InvalidArgumentWithMsg(dummyErrMsg)
	_, ok = err.(InvalidArgument)
	test.IsTrue(t, ok, "InvalidArgument interface")
	_, ok = err.(invalidArgument)
	test.IsTrue(t, ok, "invalidArgument struct")
	msg := err.Error()
	test.DeepEqual(t, dummyErrMsg, msg, "error message")
}

func TestInternalWithError(t *testing.T) {
	var ok bool
	err := InternalWithError(dummyErr)
	_, ok = err.(Internal)
	test.IsTrue(t, ok, "Internal interface")
	_, ok = err.(internal)
	test.IsTrue(t, ok, "internal struct")
	msg := err.Error()
	test.DeepEqual(t, dummyErrMsg, msg, "error message")
}

func TestInternalWithMessage(t *testing.T) {
	var ok bool
	err := InternalWithMsg(dummyErrMsg)
	_, ok = err.(Internal)
	test.IsTrue(t, ok, "Internal interface")
	_, ok = err.(internal)
	test.IsTrue(t, ok, "internal struct")
	msg := err.Error()
	test.DeepEqual(t, dummyErrMsg, msg, "error message")
}

func TestUnknownWithError(t *testing.T) {
	var ok bool
	err := UnknownWithError(dummyErr)
	_, ok = err.(Unknown)
	test.IsTrue(t, ok, "Unknown interface")
	_, ok = err.(unknown)
	test.IsTrue(t, ok, "unknown struct")
	msg := err.Error()
	test.DeepEqual(t, dummyErrMsg, msg, "error message")
}

func TestUnknownWithMessage(t *testing.T) {
	var ok bool
	err := UnknownWithMsg(dummyErrMsg)
	_, ok = err.(Unknown)
	test.IsTrue(t, ok, "Unknown interface")
	_, ok = err.(unknown)
	test.IsTrue(t, ok, "unknown struct")
	msg := err.Error()
	test.DeepEqual(t, dummyErrMsg, msg, "error message")
}
