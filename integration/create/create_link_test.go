package create

import (
	"bytes"
	"context"
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/cli/client"
	"github.com/DuarteMRAlves/maestro/internal/cli/cmd/create"
	"github.com/DuarteMRAlves/maestro/internal/cli/resources"
	"github.com/DuarteMRAlves/maestro/internal/server"
	"gotest.tools/v3/assert"
	"io/ioutil"
	"net"
	"testing"
	"time"
)

// TestCreateLinkWithServer performs integration testing on the CreateLink
// command considering operations that require the server to be running.
// It runs a maestro server and then executes a create asset command with
// predetermined arguments, verifying its output.
func TestCreateLinkWithServer(t *testing.T) {
	tests := []struct {
		name        string
		serverAddr  string
		args        []string
		expectedOut string
	}{
		{
			"create a link with all arguments",
			"localhost:50051",
			[]string{
				"link-name",
				"--source-stage",
				"source-name",
				"--source-field",
				"SourceField",
				"--target-stage",
				"target-name",
				"--target-field",
				"TargetField",
			},
			"",
		},
		{
			"create a link with required arguments",
			"localhost:50051",
			[]string{
				"link-name",
				"--source-stage",
				"source-name",
				"--target-stage",
				"target-name",
			},
			"",
		},
		{
			"create a link on custom address",
			"localhost:50052",
			[]string{
				"link-name",
				"--source-stage",
				"source-name",
				"--target-stage",
				"target-name",
				"--addr",
				"localhost:50052",
			},
			"",
		},
		{
			"create a link with invalid name",
			"localhost:50051",
			[]string{
				"invalid--name",
				"--source-stage",
				"source-name",
				"--target-stage",
				"target-name",
			},
			"invalid argument: invalid name 'invalid--name'",
		},
		{
			"create a link no such source stage",
			"localhost:50051",
			[]string{
				"link-name",
				"--source-stage",
				"does-not-exist",
				"--target-stage",
				"target-name",
			},
			"not found: source stage 'does-not-exist' not found",
		},
		{
			"create a link no such target stage",
			"localhost:50051",
			[]string{
				"link-name",
				"--source-stage",
				"source-name",
				"--target-stage",
				"does-not-exist",
			},
			"not found: target stage 'does-not-exist' not found",
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
					s.GracefulStopGrpc()
					// Wait a bit before forcefully stopping every call
					time.Sleep(10 * time.Millisecond)
					s.StopGrpc()
				}()

				// Create asset before executing command
				testResources := []*resources.Resource{
					{
						Kind: "asset",
						Spec: map[string]string{"name": "asset-name"},
					},
					{
						Kind: "stage",
						Spec: map[string]string{
							"name":  "source-name",
							"asset": "asset-name",
						},
					},
					{
						Kind: "stage",
						Spec: map[string]string{
							"name":  "target-name",
							"asset": "asset-name",
						},
					},
				}
				ctx, cancel := context.WithTimeout(
					context.Background(),
					time.Second)
				defer cancel()

				for _, r := range testResources {
					assert.NilError(
						t,
						client.CreateResource(ctx, r, test.serverAddr),
						"create resource error")
				}

				// Create link
				b := bytes.NewBufferString("")
				cmd := create.NewCmdCreateLink()
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

// TestCreateLinkWithoutServer performs integration testing on the CreateStage
// command with sets of flags that do no required the server to be running.
func TestCreateLinkWithoutServer(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectedOut string
	}{
		{
			"no name",
			[]string{},
			"invalid argument: please specify a link name",
		},
		{
			"no source stage",
			[]string{"link-name", "--target-stage", "target-name"},
			"invalid argument: please specify a source stage",
		},
		{
			"no target stage",
			[]string{"link-name", "--source-stage", "source-name"},
			"invalid argument: please specify a target stage",
		},
	}
	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				b := bytes.NewBufferString("")
				cmd := create.NewCmdCreateLink()
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
