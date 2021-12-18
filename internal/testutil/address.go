package testutil

import (
	"gotest.tools/v3/assert"
	"net"
	"testing"
)

// ListenAvailablePort returns a new Listener on an empty port where the server
// can run. The client should connect to the address using the Addr method.
func ListenAvailablePort(t *testing.T) net.Listener {
	lis, err := net.Listen("tcp", "localhost:0")
	assert.NilError(t, err, "failed to listen")
	return lis
}
