package create

import (
	"bytes"
	"context"
	"fmt"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/testutil"
	"github.com/DuarteMRAlves/maestro/internal/testutil/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gotest.tools/v3/assert"
	"io/ioutil"
	"regexp"
	"testing"
)

// TestCreateOrchestrationWithServer performs testing on the CreateOrchestration
// command considering operations that require the server to be running. It runs
// a mock maestro server and then executes a create orchestration command with
// predetermined arguments, verifying its output.
func TestCreateOrchestrationWithServer(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		validateCfg func(cfg *pb.Orchestration) bool
		response    *emptypb.Empty
		err         error
		expectedOut string
	}{
		{
			name: "create a orchestration with all arguments",
			args: []string{"orchestration-name", "--link=link1,link2"},
			validateCfg: func(cfg *pb.Orchestration) bool {
				return cfg.Name == "orchestration-name" &&
					len(cfg.Links) == 2 &&
					((cfg.Links[0] == "link1" && cfg.Links[1] == "link2") ||
						(cfg.Links[0] == "link2" && cfg.Links[1] == "link1"))
			},
			response:    &emptypb.Empty{},
			err:         nil,
			expectedOut: "",
		},
		{
			name: "create a orchestration with separate links",
			args: []string{
				"orchestration-name",
				"--link",
				"link2",
				"--link",
				"link1",
			},
			validateCfg: func(cfg *pb.Orchestration) bool {
				return cfg.Name == "orchestration-name" &&
					len(cfg.Links) == 2 &&
					((cfg.Links[0] == "link1" && cfg.Links[1] == "link2") ||
						(cfg.Links[0] == "link2" && cfg.Links[1] == "link1"))
			},
			response:    &emptypb.Empty{},
			err:         nil,
			expectedOut: "",
		},
		{
			name: "create a orchestration with required arguments",
			args: []string{"orchestration-name", "--link=link1"},
			validateCfg: func(cfg *pb.Orchestration) bool {
				return cfg.Name == "orchestration-name" &&
					len(cfg.Links) == 1 &&
					cfg.Links[0] == "link1"
			},
			response:    &emptypb.Empty{},
			err:         nil,
			expectedOut: "",
		},
		{
			name: "create a orchestration with invalid name",
			args: []string{"invalid--name", "--link=link1,link2"},
			validateCfg: func(cfg *pb.Orchestration) bool {
				return cfg.Name == "invalid--name" &&
					len(cfg.Links) == 2 &&
					((cfg.Links[0] == "link1" && cfg.Links[1] == "link2") ||
						(cfg.Links[0] == "link2" && cfg.Links[1] == "link1"))
			},
			response: nil,
			err: status.Error(
				codes.InvalidArgument,
				errdefs.InvalidArgumentWithMsg(
					"invalid name 'invalid--name'").Error()),
			expectedOut: "invalid argument: invalid name 'invalid--name'",
		},
		{
			name: "create a orchestration no such link",
			args: []string{
				"orchestration-name",
				"--link=link1,does-not-exist",
			},
			validateCfg: func(cfg *pb.Orchestration) bool {
				return cfg.Name == "orchestration-name" &&
					len(cfg.Links) == 2 &&
					((cfg.Links[0] == "link1" && cfg.Links[1] == "does-not-exist") ||
						(cfg.Links[0] == "does-not-exist" && cfg.Links[1] == "link1"))
			},
			response: nil,
			err: status.Error(
				codes.NotFound,
				errdefs.NotFoundWithMsg(
					"link 'does-not-exist' not found").Error()),
			expectedOut: "not found: link 'does-not-exist' not found",
		},
	}
	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				lis := testutil.ListenAvailablePort(t)

				addr := lis.Addr().String()
				test.args = append(test.args, "--maestro", addr)

				mockServer := mock.MaestroServer{
					OrchestrationManagementServer: &mock.OrchestrationManagementServer{
						CreateOrchestrationFn: func(
							ctx context.Context,
							cfg *pb.Orchestration,
						) (*emptypb.Empty, error) {
							if !test.validateCfg(cfg) {
								return nil, fmt.Errorf(
									"validation failed with cfg %v",
									cfg)
							}
							return test.response, test.err
						},
					},
				}
				grpcServer := mockServer.GrpcServer()
				go func() {
					err := grpcServer.Serve(lis)
					assert.NilError(t, err, "grpc server error")
				}()
				defer grpcServer.Stop()

				// Create orchestration
				b := bytes.NewBufferString("")
				cmd := NewCmdCreateOrchestration()
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
