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

// TestCreateWithServer performs integration testing on the Create command
// considering operations that require the server to be running.
// It runs a maestro server and executes the command with predetermined
// arguments, verifying its output.
func TestCreateWithServer(t *testing.T) {
	tests := []struct {
		name                  string
		args                  []string
		validateAsset         func(cfg *pb.Asset) bool
		validateStage         func(cfg *pb.Stage) bool
		validateLink          func(cfg *pb.Link) bool
		validateOrchestration func(cfg *pb.Orchestration) bool
		expectedOut           string
	}{
		{
			name: "multiple resources in a single file",
			args: []string{
				"-f",
				"../../../../tests/resources/create/resources.yml",
			},
			validateAsset: func(cfg *pb.Asset) bool {
				return equalAsset(
					&pb.Asset{Name: "asset-1", Image: "image-1"},
					cfg,
				) || equalAsset(
					&pb.Asset{Name: "asset-2", Image: "image-2"},
					cfg,
				)
			},
			validateStage: func(cfg *pb.Stage) bool {
				return equalStage(
					&pb.Stage{
						Name:    "stage-1",
						Asset:   "asset-1",
						Service: "Service1",
						Method:  "Method1",
						Address: "address-1",
						Host:    "",
						Port:    0,
					},
					cfg,
				) || equalStage(
					&pb.Stage{
						Name:    "stage-2",
						Asset:   "asset-2",
						Service: "Service2",
						Method:  "Method2",
						Address: "address-2",
						Host:    "",
						Port:    0,
					},
					cfg,
				) || equalStage(
					&pb.Stage{
						Name:    "stage-3",
						Asset:   "asset-3",
						Service: "Service3",
						Method:  "Method3",
						Address: "",
						Host:    "host-3",
						Port:    33333,
					},
					cfg,
				)
			},
			validateLink: func(cfg *pb.Link) bool {
				return equalLink(
					&pb.Link{
						Name:        "link-stage-2-stage-1",
						SourceStage: "stage-2",
						SourceField: "",
						TargetStage: "stage-1",
						TargetField: "",
					},
					cfg,
				) || equalLink(
					&pb.Link{
						Name:        "link-stage-1-stage-2",
						SourceStage: "stage-1",
						SourceField: "Field1",
						TargetStage: "stage-2",
						TargetField: "Field2",
					},
					cfg,
				)
			},
			validateOrchestration: func(cfg *pb.Orchestration) bool {
				return equalOrchestration(
					&pb.Orchestration{
						Name: "orchestration-1",
						Links: []string{
							"link-stage-1-stage-2",
							"link-stage-2-stage-1",
						},
					},
					cfg,
				) || equalOrchestration(
					&pb.Orchestration{
						Name: "orchestration-2",
						Links: []string{
							"link-stage-2-stage-1",
							"link-stage-1-stage-2",
						},
					},
					cfg,
				)
			},
			expectedOut: "",
		},
		{
			name: "multiple resources in multiple files",
			args: []string{
				"-f",
				"../../../../tests/resources/create/orchestrations.yml",
				"-f",
				"../../../../tests/resources/create/stages.yml",
				"-f",
				"../../../../tests/resources/create/links.yml",
				"-f",
				"../../../../tests/resources/create/assets.yml",
			},
			validateAsset: func(cfg *pb.Asset) bool {
				return equalAsset(
					&pb.Asset{Name: "asset-4", Image: "image-4"},
					cfg,
				) || equalAsset(
					&pb.Asset{Name: "asset-5", Image: "image-5"},
					cfg,
				) || equalAsset(
					&pb.Asset{Name: "asset-6", Image: "image-6"},
					cfg,
				)
			},
			validateStage: func(cfg *pb.Stage) bool {
				return equalStage(
					&pb.Stage{
						Name:    "stage-4",
						Asset:   "asset-4",
						Service: "",
						Method:  "",
						Address: "",
						Host:    "",
						Port:    0,
					},
					cfg,
				) || equalStage(
					&pb.Stage{
						Name:    "stage-5",
						Asset:   "",
						Service: "",
						Method:  "",
						Address: "",
						Host:    "",
						Port:    0,
					},
					cfg,
				) || equalStage(
					&pb.Stage{
						Name:    "stage-6",
						Asset:   "asset-6",
						Service: "Service6",
						Method:  "Method6",
						Address: "stage-address",
						Host:    "",
						Port:    0,
					},
					cfg,
				) || equalStage(
					&pb.Stage{
						Name:    "stage-7",
						Asset:   "asset-7",
						Service: "Service7",
						Method:  "Method7",
						Address: "",
						Host:    "stage-host",
						Port:    7777,
					},
					cfg,
				)
			},
			validateLink: func(cfg *pb.Link) bool {
				return equalLink(
					&pb.Link{
						Name:        "link-stage-4-stage-5",
						SourceStage: "stage-4",
						SourceField: "",
						TargetStage: "stage-5",
						TargetField: "",
					},
					cfg,
				) || equalLink(
					&pb.Link{
						Name:        "link-stage-5-stage-6",
						SourceStage: "stage-5",
						SourceField: "",
						TargetStage: "stage-6",
						TargetField: "Field1",
					},
					cfg,
				) || equalLink(
					&pb.Link{
						Name:        "link-stage-4-stage-6",
						SourceStage: "stage-4",
						SourceField: "",
						TargetStage: "stage-6",
						TargetField: "Field2",
					},
					cfg,
				)
			},
			validateOrchestration: func(cfg *pb.Orchestration) bool {
				return equalOrchestration(
					&pb.Orchestration{
						Name: "orchestration-3",
						Links: []string{
							"link-stage-4-stage-5",
							"link-stage-5-stage-6",
						},
					},
					cfg,
				) || equalOrchestration(
					&pb.Orchestration{
						Name: "orchestration-4",
						Links: []string{
							"link-stage-5-stage-6",
							"link-stage-4-stage-5",
							"link-stage-4-stage-6",
						},
					},
					cfg,
				)
			},
			expectedOut: "",
		},
		{
			name: "asset not found",
			args: []string{
				"-f",
				"../../../../tests/resources/create/asset_not_found.yml",
			},
			validateStage: func(cfg *pb.Stage) bool {
				return equalStage(
					&pb.Stage{
						Name:    "stage-unknown-asset",
						Asset:   "unknown-asset",
						Service: "Service1",
						Method:  "Method1",
					},
					cfg,
				)
			},
			expectedOut: "not found: asset 'unknown-asset' not found",
		},
		{
			name: "stage not found",
			args: []string{
				"-f",
				"../../../../tests/resources/create/stage_not_found.yml",
			},
			validateLink: func(cfg *pb.Link) bool {
				return equalLink(
					&pb.Link{
						Name:        "link-unknown-stage",
						SourceStage: "stage-1",
						SourceField: "",
						TargetStage: "unknown-stage",
						TargetField: "",
					},
					cfg,
				)
			},
			expectedOut: "not found: target stage 'unknown-stage' not found",
		},
		{
			name: "link not found",
			args: []string{
				"-f",
				"../../../../tests/resources/create/link_not_found.yml",
			},
			validateOrchestration: func(cfg *pb.Orchestration) bool {
				return equalOrchestration(
					&pb.Orchestration{
						Name:  "orchestration-unknown-link",
						Links: []string{"link-1", "unknown-link"},
					},
					cfg,
				)
			},
			expectedOut: "not found: link 'unknown-link' not found",
		},
	}
	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				var (
					createAssetFn func(
						ctx context.Context,
						cfg *pb.Asset,
					) (*emptypb.Empty, error)
					createStageFn func(
						ctx context.Context,
						cfg *pb.Stage,
					) (*emptypb.Empty, error)
					createLinkFn func(
						ctx context.Context,
						cfg *pb.Link,
					) (*emptypb.Empty, error)
					createOrchestrationFn func(
						ctx context.Context,
						cfg *pb.Orchestration,
					) (*emptypb.Empty, error)
				)

				lis := testutil.ListenAvailablePort(t)

				addr := lis.Addr().String()
				test.args = append(test.args, "--maestro", addr)

				if test.validateAsset != nil {
					createAssetFn = func(
						ctx context.Context,
						cfg *pb.Asset,
					) (*emptypb.Empty, error) {
						if !test.validateAsset(cfg) {
							return nil, fmt.Errorf(
								"asset validation failed with cfg %v",
								cfg)
						}
						return &emptypb.Empty{}, nil
					}
				}

				if test.validateStage != nil {
					createStageFn = func(
						ctx context.Context,
						cfg *pb.Stage,
					) (*emptypb.Empty, error) {
						if !test.validateStage(cfg) {
							return nil, fmt.Errorf(
								"stage validation failed with cfg %v",
								cfg)
						}
						if cfg.Name == "stage-unknown-asset" {
							return nil, status.Error(
								codes.NotFound,
								errdefs.NotFoundWithMsg(
									"asset 'unknown-asset' not found").Error())
						}
						return &emptypb.Empty{}, nil
					}
				}

				if test.validateLink != nil {
					createLinkFn = func(
						ctx context.Context,
						cfg *pb.Link,
					) (*emptypb.Empty, error) {
						if !test.validateLink(cfg) {
							return nil, fmt.Errorf(
								"link validation failed with cfg %v",
								cfg)
						}
						if cfg.Name == "link-unknown-stage" {
							return nil, status.Error(
								codes.NotFound,
								errdefs.NotFoundWithMsg(
									"target stage 'unknown-stage' not found").Error())
						}
						return &emptypb.Empty{}, nil
					}
				}

				if test.validateOrchestration != nil {
					createOrchestrationFn = func(
						ctx context.Context,
						cfg *pb.Orchestration,
					) (*emptypb.Empty, error) {
						if !test.validateOrchestration(cfg) {
							return nil, fmt.Errorf(
								"orchestration validation failed with cfg %v",
								cfg)
						}
						if cfg.Name == "orchestration-unknown-link" {
							return nil, status.Error(
								codes.NotFound,
								errdefs.NotFoundWithMsg(
									"link 'unknown-link' not found").Error())
						}
						return &emptypb.Empty{}, nil
					}
				}

				mockServer := mock.MaestroServer{
					AssetManagementServer: &mock.AssetManagementServer{
						CreateAssetFn: createAssetFn,
					},
					StageManagementServer: &mock.StageManagementServer{
						CreateStageFn: createStageFn,
					},
					LinkManagementServer: &mock.LinkManagementServer{
						CreateLinkFn: createLinkFn,
					},
					OrchestrationManagementServer: &mock.OrchestrationManagementServer{
						CreateOrchestrationFn: createOrchestrationFn,
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
			})
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
				"../../../../tests/resources/create/invalid_kind.yml",
			},
			"invalid argument: unknown kind: 'invalid-kind'",
		},
		{
			"invalid specs",
			[]string{
				"-f",
				"../../../../tests/resources/create/invalid_specs.yml",
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
					string(out))
				assert.NilError(t, err, "matched output")
				assert.Assert(t, matched, "output not matched")
			})
	}
}
