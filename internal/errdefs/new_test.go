package errdefs

import (
	testing2 "github.com/DuarteMRAlves/maestro/internal/testing"
	"testing"
)

func TestAlreadyExistsWithError(t *testing.T) {
	var ok bool
	err := AlreadyExistsWithError(dummyErr)
	_, ok = err.(AlreadyExists)
	testing2.IsTrue(t, ok, "AlreadyExists interface")
	_, ok = err.(alreadyExists)
	testing2.IsTrue(t, ok, "alreadyExists struct")
	msg := err.Error()
	testing2.DeepEqual(t, dummyErrMsg, msg, "error message")
}

func TestAlreadyExistsWithMsg(t *testing.T) {
	var ok bool
	err := AlreadyExistsWithMsg(dummyErrMsg)
	_, ok = err.(AlreadyExists)
	testing2.IsTrue(t, ok, "AlreadyExists interface")
	_, ok = err.(alreadyExists)
	testing2.IsTrue(t, ok, "alreadyExists struct")
	msg := err.Error()
	testing2.DeepEqual(t, dummyErrMsg, msg, "error message")
}

func TestNotFoundWithError(t *testing.T) {
	var ok bool
	err := NotFoundWithError(dummyErr)
	_, ok = err.(NotFound)
	testing2.IsTrue(t, ok, "NotFound interface")
	_, ok = err.(notFound)
	testing2.IsTrue(t, ok, "notFound struct")
	msg := err.Error()
	testing2.DeepEqual(t, dummyErrMsg, msg, "error message")
}

func TestNotFoundWithMessage(t *testing.T) {
	var ok bool
	err := NotFoundWithMsg(dummyErrMsg)
	_, ok = err.(NotFound)
	testing2.IsTrue(t, ok, "NotFound interface")
	_, ok = err.(notFound)
	testing2.IsTrue(t, ok, "notFound struct")
	msg := err.Error()
	testing2.DeepEqual(t, dummyErrMsg, msg, "error message")
}

func TestInvalidArgumentWithError(t *testing.T) {
	var ok bool
	err := InvalidArgumentWithError(dummyErr)
	_, ok = err.(InvalidArgument)
	testing2.IsTrue(t, ok, "InvalidArgument interface")
	_, ok = err.(invalidArgument)
	testing2.IsTrue(t, ok, "invalidArgument struct")
	msg := err.Error()
	testing2.DeepEqual(t, dummyErrMsg, msg, "error message")
}

func TestInvalidArgumentWithMessage(t *testing.T) {
	var ok bool
	err := InvalidArgumentWithMsg(dummyErrMsg)
	_, ok = err.(InvalidArgument)
	testing2.IsTrue(t, ok, "InvalidArgument interface")
	_, ok = err.(invalidArgument)
	testing2.IsTrue(t, ok, "invalidArgument struct")
	msg := err.Error()
	testing2.DeepEqual(t, dummyErrMsg, msg, "error message")
}

func TestInternalWithError(t *testing.T) {
	var ok bool
	err := InternalWithError(dummyErr)
	_, ok = err.(Internal)
	testing2.IsTrue(t, ok, "Internal interface")
	_, ok = err.(internal)
	testing2.IsTrue(t, ok, "internal struct")
	msg := err.Error()
	testing2.DeepEqual(t, dummyErrMsg, msg, "error message")
}

func TestInternalWithMessage(t *testing.T) {
	var ok bool
	err := InternalWithMsg(dummyErrMsg)
	_, ok = err.(Internal)
	testing2.IsTrue(t, ok, "Internal interface")
	_, ok = err.(internal)
	testing2.IsTrue(t, ok, "internal struct")
	msg := err.Error()
	testing2.DeepEqual(t, dummyErrMsg, msg, "error message")
}

func TestUnknownWithError(t *testing.T) {
	var ok bool
	err := UnknownWithError(dummyErr)
	_, ok = err.(Unknown)
	testing2.IsTrue(t, ok, "Unknown interface")
	_, ok = err.(unknown)
	testing2.IsTrue(t, ok, "unknown struct")
	msg := err.Error()
	testing2.DeepEqual(t, dummyErrMsg, msg, "error message")
}

func TestUnknownWithMessage(t *testing.T) {
	var ok bool
	err := UnknownWithMsg(dummyErrMsg)
	_, ok = err.(Unknown)
	testing2.IsTrue(t, ok, "Unknown interface")
	_, ok = err.(unknown)
	testing2.IsTrue(t, ok, "unknown struct")
	msg := err.Error()
	testing2.DeepEqual(t, dummyErrMsg, msg, "error message")
}
