package create

import (
	"bytes"
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/server"
	"gotest.tools/v3/assert"
	"io/ioutil"
	"net"
	"testing"
)

// TestCreateAssetWithServer performs integration testing on the CreateAsset
// command considering operations that require the server to be running.
// It runs a maestro server and then executes a create asset command with
// predetermined arguments, verifying its output.
func TestCreateAssetWithServer(t *testing.T) {
	tests := []struct {
		name        string
		serverAddr  string
		args        []string
		expectedOut string
	}{
		{
			"create an asset with an image",
			"localhost:50051",
			[]string{"asset-name", "--image", "image-name"},
			"",
		},
		{
			"create an asset without an image",
			"localhost:50051",
			[]string{"asset-name"},
			"",
		},
		{
			"create an asset on custom address",
			"localhost:50052",
			[]string{"asset-name", "--addr", "localhost:50052"},
			"",
		},
		{
			"create an asset invalid name",
			"localhost:50051",
			[]string{"invalid--name"},
			"invalid argument: invalid name 'invalid--name'",
		},
	}
	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				lis, err := net.Listen("tcp", test.serverAddr)
				assert.NilError(t, err, "failed to listen")
				s := server.NewBuilder().WithGrpc().Build()

				go func() {
					if err := s.ServeGrpc(lis); err != nil {
						fmt.Println(err)
						t.Errorf("Failed to serve: %v", err)
						return
					}
				}()
				defer func() {
					// Stop the server. Any calls in the test should be finished.
					// If not, an error should be raised.
					s.StopGrpc()
				}()

				b := bytes.NewBufferString("")
				cmd := NewCmdCreateAsset()
				cmd.SetOut(b)
				cmd.SetArgs(test.args)
				err = cmd.Execute()
				assert.NilError(t, err, "execute error")
				out, err := ioutil.ReadAll(b)
				assert.NilError(t, err, "read output error")
				assert.Equal(t, test.expectedOut, string(out), "output differs")
			})
	}
}

// TestCreateAssetWithoutServer performs integration testing on the CreateAsset
// command with sets of flags that do no required the server to be running.
func TestCreateAssetWithoutServer(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectedOut string
	}{
		{
			"no name",
			[]string{},
			"invalid argument: please specify the asset name",
		},
	}
	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				b := bytes.NewBufferString("")
				cmd := NewCmdCreateAsset()
				cmd.SetOut(b)
				cmd.SetArgs(test.args)
				err := cmd.Execute()
				assert.NilError(t, err, "execute error")
				out, err := ioutil.ReadAll(b)
				assert.NilError(t, err, "read output error")
				assert.Equal(t, test.expectedOut, string(out), "output differs")
			})
	}
}
