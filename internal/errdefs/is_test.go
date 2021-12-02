package errdefs

import (
	"gotest.tools/v3/assert"
	"testing"
)

func TestIsAlreadyExists(t *testing.T) {
	err := alreadyExists{dummyErr}
	assert.Assert(t, IsAlreadyExists(err), "type assertion")
}

func TestIsNotFound(t *testing.T) {
	err := notFound{dummyErr}
	assert.Assert(t, IsNotFound(err), "type assertion")
}

func TestIsInvalidArgument(t *testing.T) {
	err := invalidArgument{dummyErr}
	assert.Assert(t, IsInvalidArgument(err), "type assertion")
}

func TestIsInternal(t *testing.T) {
	err := internal{dummyErr}
	assert.Assert(t, IsInternal(err), "type assertion")
}

func TestIsUnknown(t *testing.T) {
	err := unknown{dummyErr}
	assert.Assert(t, IsUnknown(err), "type assertion")
}
