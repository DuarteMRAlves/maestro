package errdefs

import (
	testing2 "github.com/DuarteMRAlves/maestro/internal/testing"
	"testing"
)

func TestIsAlreadyExists(t *testing.T) {
	err := alreadyExists{dummyErr}
	testing2.IsTrue(t, IsAlreadyExists(err), "type assertion")
}

func TestIsNotFound(t *testing.T) {
	err := notFound{dummyErr}
	testing2.IsTrue(t, IsNotFound(err), "type assertion")
}

func TestIsInvalidArgument(t *testing.T) {
	err := invalidArgument{dummyErr}
	testing2.IsTrue(t, IsInvalidArgument(err), "type assertion")
}

func TestIsInternal(t *testing.T) {
	err := internal{dummyErr}
	testing2.IsTrue(t, IsInternal(err), "type assertion")
}

func TestIsUnknown(t *testing.T) {
	err := unknown{dummyErr}
	testing2.IsTrue(t, IsUnknown(err), "type assertion")
}
