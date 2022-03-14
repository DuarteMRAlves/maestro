package errdefs

import (
	"gotest.tools/v3/assert"
	"testing"
)

func TestIsInternal(t *testing.T) {
	err := internal{dummyErr}
	assert.Assert(t, IsInternal(err), "type assertion")
}
