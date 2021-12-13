package testutil

import (
	"gotest.tools/v3/assert"
	"net"
	"sync"
	"testing"
)

var defaultAddrLock = sync.Mutex{}

// ListenAvailablePort returns a new Listener on an empty port where the server
// can run. The client should connect to the address using the Addr method.
func ListenAvailablePort(t *testing.T) net.Listener {
	lis, err := net.Listen("tcp", "localhost:0")
	assert.NilError(t, err, "failed to listen")
	return lis
}

// LockAndListenDefaultAddr offers a synchronization mechanism for tests that
// want to test the server default address by using a global lock on the
// address. Tests should use this as little as possible in order to allow
// parallel execution. This function returns a Listener for the default address
// and locks it.
func LockAndListenDefaultAddr(t *testing.T) net.Listener {
	defaultAddrLock.Lock()
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		defaultAddrLock.Unlock()
		t.Fatalf("failed to listen on default address: %v", err)
	}
	return lis
}

// UnlockDefaultAddr offers a synchronization mechanism for tests that want to
// test the server default address by using a global lock on the address. Tests
// should use this as little as possible in order to allow parallel execution.
// This function unlocks the address and should be used when it is no longer
// required.
func UnlockDefaultAddr() {
	defaultAddrLock.Unlock()
}
