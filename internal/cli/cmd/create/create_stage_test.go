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

// TestCreateStageWithServer performs testing on the CreateStage command
// considering operations that require the server to be running. It runs a mock
// maestro server and then executes a create link command with predetermined
// arguments, verifying its output.
func TestCreateStageWithServer(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		validateCfg func(cfg *pb.Stage) bool
		response    *emptypb.Empty
		err         error
		expectedOut string
	}{
		{
			name: "create a stage with all arguments",
			args: []string{
				"stage-name",
				"--asset",
				"asset-name",
				"--service",
				"ServiceName",
				"--method",
				"MethodName",
				"--address",
				"some-address",
			},
			validateCfg: func(cfg *pb.Stage) bool {
				return cfg.Name == "stage-name" &&
					cfg.Asset == "asset-name" &&
					cfg.Service == "ServiceName" &&
					cfg.Method == "MethodName" &&
					cfg.Address == "some-address"
			},
			response:    &emptypb.Empty{},
			err:         nil,
			expectedOut: "",
		},
		{
			name: "create a stage with required arguments",
			args: []string{"stage-name"},
			validateCfg: func(cfg *pb.Stage) bool {
				return cfg.Name == "stage-name" &&
					cfg.Asset == "" &&
					cfg.Service == "" &&
					cfg.Method == "" &&
					cfg.Address == ""
			},
			response:    &emptypb.Empty{},
			err:         nil,
			expectedOut: "",
		},
		{
			name: "create a stage with invalid name",
			args: []string{"invalid--name"},
			validateCfg: func(cfg *pb.Stage) bool {
				return cfg.Name == "invalid--name" &&
					cfg.Asset == "" &&
					cfg.Service == "" &&
					cfg.Method == "" &&
					cfg.Address == ""
			},
			response: nil,
			err: status.Error(
				codes.InvalidArgument,
				errdefs.InvalidArgumentWithMsg(
					"invalid name 'invalid--name'").Error()),
			expectedOut: "invalid argument: invalid name 'invalid--name'",
		},
		{
			name: "create a stage no such asset",
			args: []string{"stage-name", "--asset", "does-not-exist"},
			validateCfg: func(cfg *pb.Stage) bool {
				return cfg.Name == "stage-name" &&
					cfg.Asset == "does-not-exist" &&
					cfg.Service == "" &&
					cfg.Method == "" &&
					cfg.Address == ""
			},
			response: nil,
			err: status.Error(
				codes.NotFound,
				errdefs.NotFoundWithMsg(
					"asset 'does-not-exist' not found").Error()),
			expectedOut: "not found: asset 'does-not-exist' not found",
		},
	}
	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				lis := testutil.ListenAvailablePort(t)

				addr := lis.Addr().String()
				test.args = append(test.args, "--addr", addr)

				mockServer := mock.MaestroServer{
					StageManagementServer: &mock.StageManagementServer{
						CreateStageFn: func(
							ctx context.Context,
							cfg *pb.Stage,
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

				b := bytes.NewBufferString("")
				cmd := NewCmdCreateStage()
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
