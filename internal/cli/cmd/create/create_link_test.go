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

// TestCreateLinkWithServer performs testing on the CreateLink command
// considering operations that require the server to be running. It runs a mock
// maestro server and then executes a create link command with predetermined
// arguments, verifying its output.
func TestCreateLinkWithServer(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		validateCfg func(cfg *pb.Link) bool
		response    *emptypb.Empty
		err         error
		expectedOut string
	}{
		{
			name: "create a link with all arguments",
			args: []string{
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
			validateCfg: func(cfg *pb.Link) bool {
				return cfg.Name == "link-name" &&
					cfg.SourceStage == "source-name" &&
					cfg.SourceField == "SourceField" &&
					cfg.TargetStage == "target-name" &&
					cfg.TargetField == "TargetField"
			},
			response:    &emptypb.Empty{},
			err:         nil,
			expectedOut: "",
		},
		{
			name: "create a link with required arguments",
			args: []string{
				"link-name",
				"--source-stage",
				"source-name",
				"--target-stage",
				"target-name",
			},
			validateCfg: func(cfg *pb.Link) bool {
				return cfg.Name == "link-name" &&
					cfg.SourceStage == "source-name" &&
					cfg.SourceField == "" &&
					cfg.TargetStage == "target-name" &&
					cfg.TargetField == ""
			},
			response:    &emptypb.Empty{},
			err:         nil,
			expectedOut: "",
		},
		{
			name: "create a link with invalid name",
			args: []string{
				"invalid--name",
				"--source-stage",
				"source-name",
				"--target-stage",
				"target-name",
			},
			validateCfg: func(cfg *pb.Link) bool {
				return cfg.Name == "invalid--name" &&
					cfg.SourceStage == "source-name" &&
					cfg.SourceField == "" &&
					cfg.TargetStage == "target-name" &&
					cfg.TargetField == ""
			},
			response: nil,
			err: status.Error(
				codes.InvalidArgument,
				errdefs.InvalidArgumentWithMsg(
					"invalid name 'invalid--name'").Error()),
			expectedOut: "invalid argument: invalid name 'invalid--name'",
		},
		{
			name: "create a link no such source stage",
			args: []string{
				"link-name",
				"--source-stage",
				"does-not-exist",
				"--target-stage",
				"target-name",
			},
			validateCfg: func(cfg *pb.Link) bool {
				return cfg.Name == "link-name" &&
					cfg.SourceStage == "does-not-exist" &&
					cfg.SourceField == "" &&
					cfg.TargetStage == "target-name" &&
					cfg.TargetField == ""
			},
			response: nil,
			err: status.Error(
				codes.NotFound,
				errdefs.NotFoundWithMsg(
					"source stage 'does-not-exist' not found").Error()),
			expectedOut: "not found: source stage 'does-not-exist' not found",
		},
		{
			name: "create a link no such target stage",
			args: []string{
				"link-name",
				"--source-stage",
				"source-name",
				"--target-stage",
				"does-not-exist",
			},
			validateCfg: func(cfg *pb.Link) bool {
				return cfg.Name == "link-name" &&
					cfg.SourceStage == "source-name" &&
					cfg.SourceField == "" &&
					cfg.TargetStage == "does-not-exist" &&
					cfg.TargetField == ""
			},
			err: status.Error(
				codes.NotFound,
				errdefs.NotFoundWithMsg(
					"target stage 'does-not-exist' not found").Error()),
			expectedOut: "not found: target stage 'does-not-exist' not found",
		},
	}
	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				lis := testutil.ListenAvailablePort(t)

				addr := lis.Addr().String()
				test.args = append(test.args, "--maestro", addr)

				mockServer := mock.MaestroServer{
					LinkManagementServer: &mock.LinkManagementServer{
						CreateLinkFn: func(
							ctx context.Context,
							cfg *pb.Link,
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

				// Create link
				b := bytes.NewBufferString("")
				cmd := NewCmdCreateLink()
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
