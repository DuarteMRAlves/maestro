package create

import (
	"bytes"
	"context"
	"fmt"
	"github.com/DuarteMRAlves/maestro/api/pb"
	ipb "github.com/DuarteMRAlves/maestro/internal/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/testutil"
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
		validateReq func(*pb.CreateStageRequest) bool
		response    *emptypb.Empty
		err         error
		expectedOut string
	}{
		{
			name: "create a stage with all arguments and address",
			args: []string{
				"stage-name",
				"--asset",
				"asset-name",
				"--service",
				"ServiceName",
				"--rpc",
				"RpcName",
				"--address",
				"some-address",
			},
			validateReq: func(req *pb.CreateStageRequest) bool {
				return req.Name == "stage-name" &&
					req.Asset == "asset-name" &&
					req.Service == "ServiceName" &&
					req.Rpc == "RpcName" &&
					req.Address == "some-address" &&
					req.Host == "" &&
					req.Port == 0
			},
			response:    &emptypb.Empty{},
			err:         nil,
			expectedOut: "",
		},
		{
			name: "create a stage with all arguments and host and port",
			args: []string{
				"stage-name",
				"--asset",
				"asset-name",
				"--service",
				"ServiceName",
				"--rpc",
				"RpcName",
				"--host",
				"some-host",
				"--port",
				"12345",
			},
			validateReq: func(req *pb.CreateStageRequest) bool {
				return req.Name == "stage-name" &&
					req.Asset == "asset-name" &&
					req.Service == "ServiceName" &&
					req.Rpc == "RpcName" &&
					req.Address == "" &&
					req.Host == "some-host" &&
					req.Port == 12345
			},
			response:    &emptypb.Empty{},
			err:         nil,
			expectedOut: "",
		},
		{
			name: "create a stage with required arguments",
			args: []string{"stage-name"},
			validateReq: func(req *pb.CreateStageRequest) bool {
				return req.Name == "stage-name" &&
					req.Asset == "" &&
					req.Service == "" &&
					req.Rpc == "" &&
					req.Address == "" &&
					req.Host == "" &&
					req.Port == 0
			},
			response:    &emptypb.Empty{},
			err:         nil,
			expectedOut: "",
		},
		{
			name: "create a stage with invalid name",
			args: []string{"invalid--name"},
			validateReq: func(req *pb.CreateStageRequest) bool {
				return req.Name == "invalid--name" &&
					req.Asset == "" &&
					req.Service == "" &&
					req.Rpc == "" &&
					req.Address == "" &&
					req.Host == "" &&
					req.Port == 0
			},
			response: nil,
			err: status.Error(
				codes.InvalidArgument,
				errdefs.InvalidArgumentWithMsg(
					"invalid name 'invalid--name'",
				).Error(),
			),
			expectedOut: "invalid argument: invalid name 'invalid--name'",
		},
		{
			name: "create a stage no such asset",
			args: []string{"stage-name", "--asset", "does-not-exist"},
			validateReq: func(req *pb.CreateStageRequest) bool {
				return req.Name == "stage-name" &&
					req.Asset == "does-not-exist" &&
					req.Service == "" &&
					req.Rpc == "" &&
					req.Address == "" &&
					req.Host == "" &&
					req.Port == 0
			},
			response: nil,
			err: status.Error(
				codes.NotFound,
				errdefs.NotFoundWithMsg(
					"asset 'does-not-exist' not found",
				).Error(),
			),
			expectedOut: "not found: asset 'does-not-exist' not found",
		},
	}
	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				lis := testutil.ListenAvailablePort(t)

				addr := lis.Addr().String()
				test.args = append(test.args, "--maestro", addr)

				mockServer := ipb.MockMaestroServer{
					StageManagementServer: &ipb.MockStageManagementServer{
						CreateStageFn: func(
							ctx context.Context,
							req *pb.CreateStageRequest,
						) (*emptypb.Empty, error) {
							if !test.validateReq(req) {
								return nil, fmt.Errorf(
									"validation failed with req %v",
									req,
								)
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
			},
		)
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
			name:        "no name",
			args:        []string{},
			expectedOut: "invalid argument: please specify a stage name",
		},
		{
			name: "address and host specified",
			args: []string{
				"stage-name",
				"--address",
				"address",
				"--host",
				"host",
			},
			expectedOut: "invalid argument: address and host options are incompatible",
		},
		{
			name: "address and port specified",
			args: []string{
				"stage-name",
				"--address",
				"address",
				"--port",
				"12345",
			},
			expectedOut: "invalid argument: address and port options are incompatible",
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
					string(out),
				)
				assert.NilError(t, err, "matched output")
				assert.Assert(t, matched, "output not matched")
			},
		)
	}
}
