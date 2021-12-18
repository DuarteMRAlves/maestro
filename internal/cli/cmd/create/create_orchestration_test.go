package create

import (
	"bytes"
	"context"
	"fmt"
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

// TestCreateOrchestrationWithServer performs integration testing on the
// CreateOrchestration command considering operations that require the server to be
// running. It runs a maestro server and then executes a create orchestration
// command with predetermined arguments, verifying its output.
func TestCreateOrchestrationWithServer(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectedOut string
	}{
		{
			name:        "create a orchestration with all arguments",
			args:        []string{"orchestration-name", "--link=link1,link2"},
			expectedOut: "",
		},
		{
			name: "create a orchestration with separate links",
			args: []string{
				"orchestration-name",
				"--link",
				"link1",
				"--link",
				"link2",
			},
			expectedOut: "",
		},
		{
			name:        "create a orchestration with required arguments",
			args:        []string{"orchestration-name", "--link=link1"},
			expectedOut: "",
		},
		{
			name:        "create a orchestration with invalid name",
			args:        []string{"invalid--name", "--link=link1,link2"},
			expectedOut: "invalid argument: invalid name 'invalid--name'",
		},
		{
			name: "create a orchestration no such link",
			args: []string{
				"orchestration-name",
				"--link=link1,does-not-exist",
			},
			expectedOut: "not found: link 'does-not-exist' not found",
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

				lis = testutil.ListenAvailablePort(t)

				addr = lis.Addr().String()

				test.args = append(test.args, "--addr", addr)

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
					{
						Kind: "link",
						Spec: &resources.LinkSpec{
							Name:        "link1",
							SourceStage: "source-name",
							TargetStage: "target-name",
						},
					},
					{
						Kind: "link",
						Spec: &resources.LinkSpec{
							Name:        "link2",
							SourceStage: "source-name",
							TargetStage: "target-name",
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

				// Create orchestration
				b := bytes.NewBufferString("")
				cmd := NewCmdCreateOrchestration()
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

// TestCreateOrchestrationWithoutServer performs integration testing on the
// CreateLink command with sets of flags that do no required the server to be
// running.
func TestCreateOrchestrationWithoutServer(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectedOut string
	}{
		{
			"no name",
			[]string{"--link=link1,link2"},
			"invalid argument: please specify a orchestration name",
		},
		{
			"no link",
			[]string{"orchestration-name"},
			"invalid argument: please specify at least one link",
		},
	}
	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				b := bytes.NewBufferString("")
				cmd := NewCmdCreateOrchestration()
				cmd.SetOut(b)
				cmd.SetArgs(test.args)
				err := cmd.Execute()
				assert.NilError(t, err, "execute error")
				out, err := ioutil.ReadAll(b)
				assert.NilError(t, err, "read output error")
				// This is not ideal but its to match the not connected error
				// with no ip. Detailed in GitHub issue
				// https://github.com/DuarteMRAlves/maestro/issues/29.
				fmt.Println(string(out))
				matched, err := regexp.MatchString(
					test.expectedOut,
					string(out))
				assert.NilError(t, err, "matched output")
				assert.Assert(t, matched, "output not matched")
			})
	}
}
