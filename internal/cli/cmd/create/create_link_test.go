package create

import (
	"bytes"
	"context"
	"github.com/DuarteMRAlves/maestro/internal/cli/client"
	"github.com/DuarteMRAlves/maestro/internal/cli/resources"
	"github.com/DuarteMRAlves/maestro/internal/server"
	"github.com/DuarteMRAlves/maestro/internal/testutil"
	"google.golang.org/grpc"
	"gotest.tools/v3/assert"
	"io/ioutil"
	"net"
	"regexp"
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
		defaultAddr bool
		args        []string
		expectedOut string
	}{
		{
			"create a link with all arguments",
			false,
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
			false,
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
			"create a link on default address",
			true,
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
			"create a link with invalid name",
			false,
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
			false,
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
			false,
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
				var (
					lis  net.Listener
					addr string
					err  error
				)

				if test.defaultAddr {
					lis = testutil.LockAndListenDefaultAddr(t)
					defer testutil.UnlockDefaultAddr()
				} else {
					lis = testutil.ListenAvailablePort(t)
				}

				addr = lis.Addr().String()

				if !test.defaultAddr {
					test.args = append(test.args, "--addr", addr)
				}
				s, err := server.NewBuilder().
					WithGrpc().
					WithLogger(testutil.NewLogger(t)).
					Build()
				assert.NilError(t, err, "build server")

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
				testResources := []*resources.Resource{
					{
						Kind: "asset",
						Spec: &resources.AssetSpec{Name: "asset-name"},
					},
					{
						Kind: "stage",
						Spec: &resources.StageSpec{
							Name: "source-name",
						},
					},
					{
						Kind: "stage",
						Spec: &resources.StageSpec{
							Name: "target-name",
						},
					},
				}
				conn, err := grpc.Dial(addr, grpc.WithInsecure())
				assert.NilError(t, err, "dial error")
				defer conn.Close()

				c := client.New(conn)

				ctx, cancel := context.WithTimeout(
					context.Background(),
					time.Second)
				defer cancel()

				for _, r := range testResources {
					assert.NilError(
						t,
						c.CreateResource(ctx, r),
						"create resource error")
				}

				// Create link
				b := bytes.NewBufferString("")
				cmd := NewCmdCreateLink()
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
		{
			"server not connected",
			[]string{
				"link-name",
				"--source-stage",
				"source-name",
				"--target-stage",
				"target-name",
			},
			`unavailable: connection error: desc = "transport: Error while dialing dial tcp .+:50051: connect: connection refused"`,
		},
	}
	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				b := bytes.NewBufferString("")
				cmd := NewCmdCreateLink()
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
