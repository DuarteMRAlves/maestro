package errdefs

import (
	"gotest.tools/v3/assert"
	"testing"
)

func TestIsInvalidArgument(t *testing.T) {
	err := invalidArgument{dummyErr}
	assert.Assert(t, IsInvalidArgument(err), "type assertion")
}

func TestIsInternal(t *testing.T) {
	err := internal{dummyErr}
	assert.Assert(t, IsInternal(err), "type assertion")
}
