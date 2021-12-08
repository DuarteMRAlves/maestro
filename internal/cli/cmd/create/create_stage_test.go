package create

import (
	"bytes"
	"context"
	"github.com/DuarteMRAlves/maestro/internal/cli/client"
	"github.com/DuarteMRAlves/maestro/internal/cli/resources"
	"github.com/DuarteMRAlves/maestro/internal/server"
	"gotest.tools/v3/assert"
	"io/ioutil"
	"net"
	"regexp"
	"testing"
	"time"
)

// TestCreateStageWithServer performs integration testing on the CreateStage
// command considering operations that require the server to be running.
// It runs a maestro server and then executes a create asset command with
// predetermined arguments, verifying its output.
func TestCreateStageWithServer(t *testing.T) {
	tests := []struct {
		name        string
		serverAddr  string
		args        []string
		expectedOut string
	}{
		{
			"create a stage with all arguments",
			"localhost:50051",
			[]string{
				"stage-name",
				"--asset",
				"asset-name",
				"--service",
				"ServiceName",
				"--method",
				"MethodName",
			},
			"",
		},
		{
			"create a stage with required arguments",
			"localhost:50051",
			[]string{"stage-name", "--asset", "asset-name"},
			"",
		},
		{
			"create an stage on custom address",
			"localhost:50052",
			[]string{
				"stage-name",
				"--asset",
				"asset-name",
				"--addr",
				"localhost:50052",
			},
			"",
		},
		{
			"create a stage with invalid name",
			"localhost:50051",
			[]string{"invalid--name", "--asset", "asset-name"},
			"invalid argument: invalid name 'invalid--name'",
		},
		{
			"create a stage no such asset",
			"localhost:50051",
			[]string{"stage-name", "--asset", "does-not-exist"},
			"not found: asset 'does-not-exist' not found",
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
						t.Errorf("Failed to serve: %v", err)
						return
					}
				}()
				defer func() {
					// Stop the server. Any calls in the test should be finished.
					// If not, an error should be raised.
					s.StopGrpc()
				}()

				// Create asset before executing command
				ctx, cancel := context.WithTimeout(
					context.Background(),
					time.Second)
				defer cancel()

				assert.NilError(
					t,
					client.CreateAsset(
						ctx,
						&resources.AssetResource{Name: "asset-name"},
						test.serverAddr),
					"create asset error")

				b := bytes.NewBufferString("")
				cmd := NewCmdCreateStage()
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

// TestCreateStageWithoutServer performs integration testing on the CreateStage
// command with sets of flags that do no required the server to be running.
func TestCreateStageWithoutServer(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectedOut string
	}{
		{
			"no name",
			[]string{},
			"invalid argument: please specify a stage name",
		},
		{
			"no asset",
			[]string{"stage-name"},
			"invalid argument: please specify an asset",
		},
		{
			"server not connected",
			[]string{"stage-name", "--asset", "asset-name"},
			`unavailable: connection error: desc = "transport: Error while dialing dial tcp .+:50051: connect: connection refused"`,
		},
	}
	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				b := bytes.NewBufferString("")
				cmd := NewCmdCreateStage()
				cmd.SetOut(b)
				cmd.SetArgs(test.args)
				err := cmd.Execute()
				assert.NilError(t, err, "execute error")
				out, err := ioutil.ReadAll(b)
				assert.NilError(t, err, "read output error")
				// This is not ideal but its to match the not connected error
				// with no ip. Detailed in GitHub issue
				// https://github.com/DuarteMRAlves/maestro/issues/29.
				matched, err := regexp.MatchString(
					test.expectedOut,
					string(out))
				assert.NilError(t, err, "matched output")
				assert.Assert(t, matched, "output not matched")
			})
	}
}
