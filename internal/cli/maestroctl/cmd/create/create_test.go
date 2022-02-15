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

// TestCreateWithServer performs integration testing on the Create command
// considering operations that require the server to be running.
// It runs a maestro server and executes the command with predetermined
// arguments, verifying its output.
func TestCreateWithServer(t *testing.T) {
	tests := []struct {
		name                  string
		args                  []string
		validateAsset         func(*pb.CreateAssetRequest) bool
		validateStage         func(*pb.CreateStageRequest) bool
		validateLink          func(*pb.CreateLinkRequest) bool
		validateOrchestration func(*pb.CreateOrchestrationRequest) bool
		validateExecution     func(request *pb.StartExecutionRequest) bool
		expectedOut           string
	}{
		{
			name: "multiple resources in a single file",
			args: []string{
				"-f",
				"../../../../../tests/resources/create/resources.yml",
			},
			validateAsset: func(req *pb.CreateAssetRequest) bool {
				return equalCreateAssetRequest(
					&pb.CreateAssetRequest{Name: "asset-1", Image: "image-1"},
					req,
				) || equalCreateAssetRequest(
					&pb.CreateAssetRequest{Name: "asset-2", Image: "image-2"},
					req,
				)
			},
			validateStage: func(req *pb.CreateStageRequest) bool {
				return equalCreateStageRequest(
					&pb.CreateStageRequest{
						Name:    "stage-1",
						Asset:   "asset-1",
						Service: "Service1",
						Rpc:     "Rpc1",
						Address: "address-1",
						Host:    "",
						Port:    0,
					},
					req,
				) || equalCreateStageRequest(
					&pb.CreateStageRequest{
						Name:    "stage-2",
						Asset:   "asset-2",
						Service: "Service2",
						Rpc:     "Rpc2",
						Address: "address-2",
						Host:    "",
						Port:    0,
					},
					req,
				) || equalCreateStageRequest(
					&pb.CreateStageRequest{
						Name:    "stage-3",
						Asset:   "asset-3",
						Service: "Service3",
						Rpc:     "Rpc3",
						Address: "",
						Host:    "host-3",
						Port:    33333,
					},
					req,
				)
			},
			validateLink: func(req *pb.CreateLinkRequest) bool {
				return equalCreateLinkRequest(
					&pb.CreateLinkRequest{
						Name:        "link-stage-2-stage-1",
						SourceStage: "stage-2",
						SourceField: "",
						TargetStage: "stage-1",
						TargetField: "",
					},
					req,
				) || equalCreateLinkRequest(
					&pb.CreateLinkRequest{
						Name:        "link-stage-1-stage-2",
						SourceStage: "stage-1",
						SourceField: "Field1",
						TargetStage: "stage-2",
						TargetField: "Field2",
					},
					req,
				)
			},
			validateOrchestration: func(req *pb.CreateOrchestrationRequest) bool {
				return equalCreateOrchestrationRequest(
					&pb.CreateOrchestrationRequest{
						Name: "orchestration-1",
					},
					req,
				) || equalCreateOrchestrationRequest(
					&pb.CreateOrchestrationRequest{
						Name: "orchestration-2",
					},
					req,
				)
			},
			validateExecution: func(req *pb.StartExecutionRequest) bool {
				return equalStartExecutionRequest(
					&pb.StartExecutionRequest{
						Orchestration: "orchestration-1",
					},
					req,
				) || equalStartExecutionRequest(
					&pb.StartExecutionRequest{
						Orchestration: "orchestration-2",
					},
					req,
				)
			},
			expectedOut: "",
		},
		{
			name: "multiple resources in multiple files",
			args: []string{
				"-f",
				"../../../../../tests/resources/create/orchestrations.yml",
				"-f",
				"../../../../../tests/resources/create/stages.yml",
				"-f",
				"../../../../../tests/resources/create/links.yml",
				"-f",
				"../../../../../tests/resources/create/assets.yml",
			},
			validateAsset: func(req *pb.CreateAssetRequest) bool {
				return equalCreateAssetRequest(
					&pb.CreateAssetRequest{Name: "asset-4", Image: "image-4"},
					req,
				) || equalCreateAssetRequest(
					&pb.CreateAssetRequest{Name: "asset-5", Image: "image-5"},
					req,
				) || equalCreateAssetRequest(
					&pb.CreateAssetRequest{Name: "asset-6", Image: "image-6"},
					req,
				)
			},
			validateStage: func(req *pb.CreateStageRequest) bool {
				return equalCreateStageRequest(
					&pb.CreateStageRequest{
						Name:    "stage-4",
						Asset:   "asset-4",
						Service: "",
						Rpc:     "",
						Address: "",
						Host:    "",
						Port:    0,
					},
					req,
				) || equalCreateStageRequest(
					&pb.CreateStageRequest{
						Name:    "stage-5",
						Asset:   "",
						Service: "",
						Rpc:     "",
						Address: "",
						Host:    "",
						Port:    0,
					},
					req,
				) || equalCreateStageRequest(
					&pb.CreateStageRequest{
						Name:    "stage-6",
						Asset:   "asset-6",
						Service: "Service6",
						Rpc:     "Rpc6",
						Address: "stage-address",
						Host:    "",
						Port:    0,
					},
					req,
				) || equalCreateStageRequest(
					&pb.CreateStageRequest{
						Name:    "stage-7",
						Asset:   "asset-7",
						Service: "Service7",
						Rpc:     "Rpc7",
						Address: "",
						Host:    "stage-host",
						Port:    7777,
					},
					req,
				)
			},
			validateLink: func(req *pb.CreateLinkRequest) bool {
				return equalCreateLinkRequest(
					&pb.CreateLinkRequest{
						Name:        "link-stage-4-stage-5",
						SourceStage: "stage-4",
						SourceField: "",
						TargetStage: "stage-5",
						TargetField: "",
					},
					req,
				) || equalCreateLinkRequest(
					&pb.CreateLinkRequest{
						Name:        "link-stage-5-stage-6",
						SourceStage: "stage-5",
						SourceField: "",
						TargetStage: "stage-6",
						TargetField: "Field1",
					},
					req,
				) || equalCreateLinkRequest(
					&pb.CreateLinkRequest{
						Name:        "link-stage-4-stage-6",
						SourceStage: "stage-4",
						SourceField: "",
						TargetStage: "stage-6",
						TargetField: "Field2",
					},
					req,
				)
			},
			validateOrchestration: func(req *pb.CreateOrchestrationRequest) bool {
				return equalCreateOrchestrationRequest(
					&pb.CreateOrchestrationRequest{
						Name: "orchestration-3",
					},
					req,
				) || equalCreateOrchestrationRequest(
					&pb.CreateOrchestrationRequest{
						Name: "orchestration-4",
					},
					req,
				)
			},
			validateExecution: func(req *pb.StartExecutionRequest) bool {
				return equalStartExecutionRequest(
					&pb.StartExecutionRequest{
						Orchestration: "orchestration-3",
					},
					req,
				) || equalStartExecutionRequest(
					&pb.StartExecutionRequest{
						Orchestration: "orchestration-4",
					},
					req,
				)
			},
			expectedOut: "",
		},
		{
			name: "filter resources",
			args: []string{
				"-f",
				"../../../../../tests/resources/create/resources.yml",
				"asset-2",
				"stage-1",
				"stage-2",
				"link-stage-2-stage-1",
				"orchestration-2",
			},
			validateAsset: func(req *pb.CreateAssetRequest) bool {
				return equalCreateAssetRequest(
					&pb.CreateAssetRequest{Name: "asset-2", Image: "image-2"},
					req,
				)
			},
			validateStage: func(req *pb.CreateStageRequest) bool {
				return equalCreateStageRequest(
					&pb.CreateStageRequest{
						Name:    "stage-1",
						Asset:   "asset-1",
						Service: "Service1",
						Rpc:     "Rpc1",
						Address: "address-1",
						Host:    "",
						Port:    0,
					},
					req,
				) || equalCreateStageRequest(
					&pb.CreateStageRequest{
						Name:    "stage-2",
						Asset:   "asset-2",
						Service: "Service2",
						Rpc:     "Rpc2",
						Address: "address-2",
						Host:    "",
						Port:    0,
					},
					req,
				)
			},
			validateLink: func(req *pb.CreateLinkRequest) bool {
				return equalCreateLinkRequest(
					&pb.CreateLinkRequest{
						Name:        "link-stage-2-stage-1",
						SourceStage: "stage-2",
						SourceField: "",
						TargetStage: "stage-1",
						TargetField: "",
					},
					req,
				)
			},
			validateOrchestration: func(req *pb.CreateOrchestrationRequest) bool {
				return equalCreateOrchestrationRequest(
					&pb.CreateOrchestrationRequest{
						Name: "orchestration-2",
					},
					req,
				)
			},
			validateExecution: func(req *pb.StartExecutionRequest) bool {
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
			name: "asset not found",
			args: []string{
				"-f",
				"../../../../../tests/resources/create/asset_not_found.yml",
			},
			validateStage: func(req *pb.CreateStageRequest) bool {
				return equalCreateStageRequest(
					&pb.CreateStageRequest{
						Name:    "stage-unknown-asset",
						Asset:   "unknown-asset",
						Service: "Service1",
						Rpc:     "Rpc1",
					},
					req,
				)
			},
			expectedOut: "not found: asset 'unknown-asset' not found",
		},
		{
			name: "stage not found",
			args: []string{
				"-f",
				"../../../../../tests/resources/create/stage_not_found.yml",
			},
			validateLink: func(req *pb.CreateLinkRequest) bool {
				return equalCreateLinkRequest(
					&pb.CreateLinkRequest{
						Name:        "link-unknown-stage",
						SourceStage: "stage-1",
						SourceField: "",
						TargetStage: "unknown-stage",
						TargetField: "",
					},
					req,
				)
			},
			expectedOut: "not found: target stage 'unknown-stage' not found",
		},
	}
	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				var (
					createAssetFn func(
						context.Context,
						*pb.CreateAssetRequest,
					) (*emptypb.Empty, error)
					createStageFn func(
						context.Context,
						*pb.CreateStageRequest,
					) (*emptypb.Empty, error)
					createLinkFn func(
						context.Context,
						*pb.CreateLinkRequest,
					) (*emptypb.Empty, error)
					createOrchestrationFn func(
						context.Context,
						*pb.CreateOrchestrationRequest,
					) (*emptypb.Empty, error)
					startExecutionFn func(
						context.Context,
						*pb.StartExecutionRequest,
					) (*emptypb.Empty, error)
				)

				lis := util.NewTestListener(t)

				addr := lis.Addr().String()
				test.args = append(test.args, "--maestro", addr)

				if test.validateAsset != nil {
					createAssetFn = func(
						ctx context.Context,
						req *pb.CreateAssetRequest,
					) (*emptypb.Empty, error) {
						if !test.validateAsset(req) {
							return nil, fmt.Errorf(
								"asset validation failed with req %v",
								req,
							)
						}
						return &emptypb.Empty{}, nil
					}
				}

				if test.validateStage != nil {
					createStageFn = func(
						ctx context.Context,
						req *pb.CreateStageRequest,
					) (*emptypb.Empty, error) {
						if !test.validateStage(req) {
							return nil, fmt.Errorf(
								"stage validation failed with req %v",
								req,
							)
						}
						if req.Name == "stage-unknown-asset" {
							return nil, status.Error(
								codes.NotFound,
								errdefs.NotFoundWithMsg(
									"asset 'unknown-asset' not found",
								).Error(),
							)
						}
						return &emptypb.Empty{}, nil
					}
				}

				if test.validateLink != nil {
					createLinkFn = func(
						ctx context.Context,
						req *pb.CreateLinkRequest,
					) (*emptypb.Empty, error) {
						if !test.validateLink(req) {
							return nil, fmt.Errorf(
								"link validation failed with req %v",
								req,
							)
						}
						if req.Name == "link-unknown-stage" {
							return nil, status.Error(
								codes.NotFound,
								errdefs.NotFoundWithMsg(
									"target stage 'unknown-stage' not found",
								).Error(),
							)
						}
						return &emptypb.Empty{}, nil
					}
				}

				if test.validateOrchestration != nil {
					createOrchestrationFn = func(
						ctx context.Context,
						req *pb.CreateOrchestrationRequest,
					) (*emptypb.Empty, error) {
						if !test.validateOrchestration(req) {
							return nil, fmt.Errorf(
								"orchestration validation failed with req %v",
								req,
							)
						}
						if req.Name == "orchestration-unknown-link" {
							return nil, status.Error(
								codes.NotFound,
								errdefs.NotFoundWithMsg(
									"link 'unknown-link' not found",
								).Error(),
							)
						}
						return &emptypb.Empty{}, nil
					}
				}

				if test.validateExecution != nil {
					startExecutionFn = func(
						ctx context.Context,
						req *pb.StartExecutionRequest,
					) (*emptypb.Empty, error) {
						if !test.validateExecution(req) {
							return nil, fmt.Errorf(
								"start execution validation failed with req %v",
								req,
							)
						}
						return &emptypb.Empty{}, nil
					}
				}

				mockServer := ipb.MockMaestroServer{
					ArchitectureManagementServer: &ipb.MockArchitectureManagementServer{
						CreateAssetFn:         createAssetFn,
						CreateOrchestrationFn: createOrchestrationFn,
						CreateStageFn:         createStageFn,
						CreateLinkFn:          createLinkFn,
					},
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
				cmd := NewCmdCreate()
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

// TestCreateWithServer performs integration testing on the Create command
// considering operations that do not require the server to be running.
func TestCreateWithoutServer(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectedOut string
	}{
		{
			"no files",
			[]string{},
			"invalid argument: please specify input files",
		},
		{
			"no such file",
			[]string{"-f", "missing_file.yml"},
			"invalid argument: open missing_file.yml: no such file or directory",
		},
		{
			"invalid kind",
			[]string{
				"-f",
				"../../../../../tests/resources/create/invalid_kind.yml",
			},
			"invalid argument: unknown kind: 'invalid-kind'",
		},
		{
			"invalid specs",
			[]string{
				"-f",
				"../../../../../tests/resources/create/invalid_specs.yml",
			},
			"invalid argument: unknown spec fields: invalid_spec_1,invalid_spec_2",
		},
	}
	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				b := bytes.NewBufferString("")
				cmd := NewCmdCreate()
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

func equalCreateAssetRequest(
	expected *pb.CreateAssetRequest,
	actual *pb.CreateAssetRequest,
) bool {
	return expected.Name == actual.Name && expected.Image == actual.Image
}

func equalCreateStageRequest(
	expected *pb.CreateStageRequest,
	actual *pb.CreateStageRequest,
) bool {
	return expected.Name == actual.Name &&
		expected.Asset == actual.Asset &&
		expected.Service == actual.Service &&
		expected.Rpc == actual.Rpc &&
		expected.Address == actual.Address &&
		expected.Host == actual.Host &&
		expected.Port == actual.Port
}

func equalCreateLinkRequest(
	expected *pb.CreateLinkRequest,
	actual *pb.CreateLinkRequest,
) bool {
	return expected.Name == actual.Name &&
		expected.SourceStage == actual.SourceStage &&
		expected.SourceField == actual.SourceField &&
		expected.TargetStage == actual.TargetStage &&
		expected.TargetField == actual.TargetField
}

func equalCreateOrchestrationRequest(
	expected *pb.CreateOrchestrationRequest,
	actual *pb.CreateOrchestrationRequest,
) bool {
	return expected.Name == actual.Name
}

func equalStartExecutionRequest(
	expected *pb.StartExecutionRequest,
	actual *pb.StartExecutionRequest,
) bool {
	return expected.Orchestration == actual.Orchestration
}