package errdefs

import (
	"github.com/DuarteMRAlves/maestro/internal/test"
	"testing"
)

func TestIsAlreadyExists(t *testing.T) {
	err := alreadyExists{dummyErr}
	test.IsTrue(t, IsAlreadyExists(err), "type assertion")
}

func TestIsNotFound(t *testing.T) {
	err := notFound{dummyErr}
	test.IsTrue(t, IsNotFound(err), "type assertion")
}

func TestIsInvalidArgument(t *testing.T) {
	err := invalidArgument{dummyErr}
	test.IsTrue(t, IsInvalidArgument(err), "type assertion")
}

func TestIsInternal(t *testing.T) {
	err := internal{dummyErr}
	test.IsTrue(t, IsInternal(err), "type assertion")
}

func TestIsUnknown(t *testing.T) {
	err := unknown{dummyErr}
	test.IsTrue(t, IsUnknown(err), "type assertion")
}
