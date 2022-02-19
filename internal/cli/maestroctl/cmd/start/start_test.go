package start

import (
	"bytes"
	"context"
	"fmt"
	"github.com/DuarteMRAlves/maestro/api/pb"
	ipb "github.com/DuarteMRAlves/maestro/internal/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gotest.tools/v3/assert"
	"io/ioutil"
	"net"
	"testing"
)

// TestCreateWithServer performs integration testing on the Create command
// considering operations that require the server to be running.
// It runs a maestro server and executes the command with predetermined
// arguments, verifying its output.
func TestCreateWithServer(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		validate    func(request *pb.StartExecutionRequest) bool
		expectedOut string
	}{
		{
			name: "default orchestration",
			args: []string{},
			validate: func(req *pb.StartExecutionRequest) bool {
				return equalStartExecutionRequest(
					&pb.StartExecutionRequest{
						Orchestration: "",
					},
					req,
				)
			},
			expectedOut: "",
		},
		{
			name: "custom orchestration",
			args: []string{"orchestration-2"},
			validate: func(req *pb.StartExecutionRequest) bool {
				return equalStartExecutionRequest(
					&pb.StartExecutionRequest{
						Orchestration: "orchestration-2",
					},
					req,
				)
			},
			expectedOut: "",
		},
		{
			name: "orchestration not found",
			args: []string{"unknown-orchestration"},
			validate: func(req *pb.StartExecutionRequest) bool {
				return equalStartExecutionRequest(
					&pb.StartExecutionRequest{
						Orchestration: "unknown-orchestration",
					},
					req,
				)
			},
			expectedOut: "not found: orchestration 'unknown-orchestration' not found",
		},
	}
	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				var (
					startExecutionFn func(
						context.Context,
						*pb.StartExecutionRequest,
					) (*emptypb.Empty, error)
				)

				lis, err := net.Listen("tcp", "localhost:0")
				assert.NilError(t, err, "failed to listen")

				addr := lis.Addr().String()
				test.args = append(test.args, "--maestro", addr)

				if test.validate != nil {
					startExecutionFn = func(
						ctx context.Context,
						req *pb.StartExecutionRequest,
					) (*emptypb.Empty, error) {
						if !test.validate(req) {
							return nil, fmt.Errorf(
								"start execution validation failed with req %v",
								req,
							)
						}
						if req.Orchestration == "unknown-orchestration" {
							return nil, status.Error(
								codes.NotFound,
								errdefs.NotFoundWithMsg(
									"orchestration 'unknown-orchestration' not found",
								).Error(),
							)
						}
						return &emptypb.Empty{}, nil
					}
				}

				mockServer := ipb.MockMaestroServer{
					ExecutionManagementServer: &ipb.MockExecutionManagementServer{
						StartExecutionFn: startExecutionFn,
					},
				}
				grpcServer := mockServer.GrpcServer()
				go func() {
					err := grpcServer.Serve(lis)
					assert.NilError(t, err, "grpc server error")
				}()
				defer grpcServer.Stop()

				b := bytes.NewBufferString("")
				cmd := NewCmdStart()
				cmd.SetOut(b)
				cmd.SetArgs(test.args)
				err = cmd.Execute()
				assert.NilError(t, err, "execute error")
				out, err := ioutil.ReadAll(b)
				assert.NilError(t, err, "read output error")
				assert.Equal(t, test.expectedOut, string(out), "output differs")
			},
		)
	}
}

func equalStartExecutionRequest(
	expected *pb.StartExecutionRequest,
	actual *pb.StartExecutionRequest,
) bool {
	return expected.Orchestration == actual.Orchestration
}
