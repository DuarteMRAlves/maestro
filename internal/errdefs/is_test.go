package errdefs

import (
	"github.com/DuarteMRAlves/maestro/internal/assert"
	"testing"
)

func TestIsAlreadyExists(t *testing.T) {
	err := alreadyExists{dummyErr}
	assert.IsTrue(t, IsAlreadyExists(err), "type assertion")
}

func TestIsNotFound(t *testing.T) {
	err := notFound{dummyErr}
	assert.IsTrue(t, IsNotFound(err), "type assertion")
}

func TestIsInvalidArgument(t *testing.T) {
	err := invalidArgument{dummyErr}
	assert.IsTrue(t, IsInvalidArgument(err), "type assertion")
}

func TestIsInternal(t *testing.T) {
	err := internal{dummyErr}
	assert.IsTrue(t, IsInternal(err), "type assertion")
}

func TestIsUnknown(t *testing.T) {
	err := unknown{dummyErr}
	assert.IsTrue(t, IsUnknown(err), "type assertion")
}
