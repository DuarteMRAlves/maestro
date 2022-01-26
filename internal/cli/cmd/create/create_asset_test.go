package create

import (
	"bytes"
	"context"
	"fmt"
	"github.com/DuarteMRAlves/maestro/api/pb"
	ipb "github.com/DuarteMRAlves/maestro/internal/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gotest.tools/v3/assert"
	"io/ioutil"
	"regexp"
	"testing"
)

// TestCreateAssetWithServer performs testing on the CreateAsset command
// considering operations that require the server to be running. It runs a mock
// maestro server and then executes a create asset command with predetermined
// arguments, verifying its output.
func TestCreateAssetWithServer(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		validateCfg func(req *pb.CreateAssetRequest) bool
		response    *emptypb.Empty
		err         error
		expectedOut string
	}{
		{
			name: "create an asset with an image",
			args: []string{"asset-name", "--image", "image-name"},
			validateCfg: func(req *pb.CreateAssetRequest) bool {
				return req.Name == "asset-name" && req.Image == "image-name"
			},
			response:    &emptypb.Empty{},
			err:         nil,
			expectedOut: "",
		},
		{
			name: "create an asset without an image",
			args: []string{"asset-name"},
			validateCfg: func(req *pb.CreateAssetRequest) bool {
				return req.Name == "asset-name" && req.Image == ""
			},
			response:    &emptypb.Empty{},
			err:         nil,
			expectedOut: "",
		},
		{
			name: "create an asset invalid name",
			args: []string{"invalid--name"},
			validateCfg: func(req *pb.CreateAssetRequest) bool {
				return req.Name == "invalid--name" && req.Image == ""
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
	}
	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				lis := util.NewTestListener(t)

				addr := lis.Addr().String()
				test.args = append(test.args, "--maestro", addr)

				mockServer := ipb.MockMaestroServer{
					AssetManagementServer: &ipb.MockAssetManagementServer{
						CreateAssetFn: func(
							ctx context.Context,
							req *pb.CreateAssetRequest,
						) (*emptypb.Empty, error) {
							if !test.validateCfg(req) {
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
				cmd := NewCmdCreateAsset()
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
