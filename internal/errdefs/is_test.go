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

func TestIsUnknown(t *testing.T) {
	err := unknown{dummyErr}
	assert.IsTrue(t, IsUnknown(err), "type assertion")
}
