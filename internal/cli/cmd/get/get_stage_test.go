package get

import (
	"bytes"
	"fmt"
	"github.com/DuarteMRAlves/maestro/api/pb"
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"github.com/DuarteMRAlves/maestro/internal/testutil"
	"github.com/DuarteMRAlves/maestro/internal/testutil/mock"
	"github.com/pterm/pterm"
	"gotest.tools/v3/assert"
	"io/ioutil"
	"testing"
)

// TestGetStage_CorrectDisplay performs testing on the GetStage command
// considering operations that produce table outputs. It runs a mock maestro
// server and then executes a get stage command with predetermined arguments,
// verifying its output by comparing with an expected table.
func TestGetStage_CorrectDisplay(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		validateQuery func(query *pb.Stage) bool
		responses     []*pb.Stage
		output        [][]string
	}{
		{
			name: "empty stages",
			args: []string{},
			validateQuery: func(query *pb.Stage) bool {
				return query.Name == "" &&
					query.Phase == "" &&
					query.Asset == "" &&
					query.Service == "" &&
					query.Rpc == "" &&
					query.Address == ""
			},
			responses: []*pb.Stage{},
			output: [][]string{
				{
					NameText,
					PhaseText,
					AssetText,
					ServiceText,
					RpcText,
					AddressText,
				},
			},
		},
		{
			name: "one stage",
			args: []string{},
			validateQuery: func(query *pb.Stage) bool {
				return query.Name == "" &&
					query.Phase == "" &&
					query.Asset == "" &&
					query.Service == "" &&
					query.Rpc == "" &&
					query.Address == ""
			},
			responses: []*pb.Stage{pbStageForNum(0, apitypes.StageRunning)},
			output: [][]string{
				{
					NameText,
					PhaseText,
					AssetText,
					ServiceText,
					RpcText,
					AddressText,
				},
				{
					stageNameForNum(0),
					string(apitypes.StageRunning),
					assetNameForNum(0),
					stageServiceForNum(0),
					stageRpcForNum(0),
					stageAddressForNum(0),
				},
			},
		},
		{
			name: "multiple stages",
			args: []string{},
			validateQuery: func(query *pb.Stage) bool {
				return query.Name == "" &&
					query.Phase == "" &&
					query.Asset == "" &&
					query.Service == "" &&
					query.Rpc == "" &&
					query.Address == ""
			},
			responses: []*pb.Stage{
				pbStageForNum(0, apitypes.StagePending),
				pbStageForNum(2, apitypes.StageRunning),
				pbStageForNum(1, apitypes.StageFailed),
			},
			output: [][]string{
				{
					NameText,
					PhaseText,
					AssetText,
					ServiceText,
					RpcText,
					AddressText,
				},
				{
					stageNameForNum(0),
					string(apitypes.StagePending),
					assetNameForNum(0),
					stageServiceForNum(0),
					stageRpcForNum(0),
					stageAddressForNum(0),
				},
				{
					stageNameForNum(1),
					string(apitypes.StageFailed),
					assetNameForNum(1),
					stageServiceForNum(1),
					stageRpcForNum(1),
					stageAddressForNum(1),
				},
				{
					stageNameForNum(2),
					string(apitypes.StageRunning),
					assetNameForNum(2),
					stageServiceForNum(2),
					stageRpcForNum(2),
					stageAddressForNum(2),
				},
			},
		},
		{
			name: "filter by name",
			args: []string{stageNameForNum(2)},
			validateQuery: func(query *pb.Stage) bool {
				return query.Name == stageNameForNum(2) &&
					query.Phase == "" &&
					query.Asset == "" &&
					query.Service == "" &&
					query.Rpc == "" &&
					query.Address == ""
			},
			responses: []*pb.Stage{
				pbStageForNum(2, apitypes.StageSucceeded),
			},
			output: [][]string{
				{
					NameText,
					PhaseText,
					AssetText,
					ServiceText,
					RpcText,
					AddressText,
				},
				{
					stageNameForNum(2),
					string(apitypes.StageSucceeded),
					assetNameForNum(2),
					stageServiceForNum(2),
					stageRpcForNum(2),
					stageAddressForNum(2),
				},
			},
		},
		{
			name: "filter by phase",
			args: []string{"--phase", string(apitypes.StageRunning)},
			validateQuery: func(query *pb.Stage) bool {
				return query.Name == "" &&
					query.Phase == string(apitypes.StageRunning) &&
					query.Asset == "" &&
					query.Service == "" &&
					query.Rpc == "" &&
					query.Address == ""
			},
			responses: []*pb.Stage{
				pbStageForNum(1, apitypes.StageRunning),
			},
			output: [][]string{
				{
					NameText,
					PhaseText,
					AssetText,
					ServiceText,
					RpcText,
					AddressText,
				},
				{
					stageNameForNum(1),
					string(apitypes.StageRunning),
					assetNameForNum(1),
					stageServiceForNum(1),
					stageRpcForNum(1),
					stageAddressForNum(1),
				},
			},
		},
		{
			name: "filter by asset",
			args: []string{"--asset", assetNameForNum(2)},
			validateQuery: func(query *pb.Stage) bool {
				return query.Name == "" &&
					query.Phase == "" &&
					query.Asset == assetNameForNum(2) &&
					query.Service == "" &&
					query.Rpc == "" &&
					query.Address == ""
			},
			responses: []*pb.Stage{
				pbStageForNum(2, apitypes.StagePending),
			},
			output: [][]string{
				{
					NameText,
					PhaseText,
					AssetText,
					ServiceText,
					RpcText,
					AddressText,
				},
				{
					stageNameForNum(2),
					string(apitypes.StagePending),
					assetNameForNum(2),
					stageServiceForNum(2),
					stageRpcForNum(2),
					stageAddressForNum(2),
				},
			},
		},
		{
			name: "filter by service",
			args: []string{"--service", stageServiceForNum(0)},
			validateQuery: func(query *pb.Stage) bool {
				return query.Name == "" &&
					query.Phase == "" &&
					query.Asset == "" &&
					query.Service == stageServiceForNum(0) &&
					query.Rpc == "" &&
					query.Address == ""
			},
			responses: []*pb.Stage{
				pbStageForNum(0, apitypes.StageRunning),
			},
			output: [][]string{
				{
					NameText,
					PhaseText,
					AssetText,
					ServiceText,
					RpcText,
					AddressText,
				},
				{
					stageNameForNum(0),
					string(apitypes.StageRunning),
					assetNameForNum(0),
					stageServiceForNum(0),
					stageRpcForNum(0),
					stageAddressForNum(0),
				},
			},
		},
		{
			name: "filter by rpc",
			args: []string{"--rpc", stageRpcForNum(1)},
			validateQuery: func(query *pb.Stage) bool {
				return query.Name == "" &&
					query.Phase == "" &&
					query.Asset == "" &&
					query.Service == "" &&
					query.Rpc == stageRpcForNum(1) &&
					query.Address == ""
			},
			responses: []*pb.Stage{
				pbStageForNum(1, apitypes.StagePending),
			},
			output: [][]string{
				{
					NameText,
					PhaseText,
					AssetText,
					ServiceText,
					RpcText,
					AddressText,
				},
				{
					stageNameForNum(1),
					string(apitypes.StagePending),
					assetNameForNum(1),
					stageServiceForNum(1),
					stageRpcForNum(1),
					stageAddressForNum(1),
				},
			},
		},
		{
			name: "no such name",
			args: []string{stageNameForNum(3)},
			validateQuery: func(query *pb.Stage) bool {
				return query.Name == stageNameForNum(3) &&
					query.Phase == "" &&
					query.Asset == "" &&
					query.Service == "" &&
					query.Rpc == "" &&
					query.Address == ""
			},
			responses: []*pb.Stage{},
			output: [][]string{
				{
					NameText,
					PhaseText,
					AssetText,
					ServiceText,
					RpcText,
					AddressText,
				},
			},
		},
		{
			name: "no such phase",
			args: []string{"--phase", string(apitypes.StagePending)},
			validateQuery: func(query *pb.Stage) bool {
				return query.Name == "" &&
					query.Phase == string(apitypes.StagePending) &&
					query.Asset == "" &&
					query.Service == "" &&
					query.Rpc == "" &&
					query.Address == ""
			},
			responses: []*pb.Stage{},
			output: [][]string{
				{
					NameText,
					PhaseText,
					AssetText,
					ServiceText,
					RpcText,
					AddressText,
				},
			},
		},
		{
			name: "no such asset",
			args: []string{"--asset", assetNameForNum(3)},
			validateQuery: func(query *pb.Stage) bool {
				return query.Name == "" &&
					query.Phase == "" &&
					query.Asset == assetNameForNum(3) &&
					query.Service == "" &&
					query.Rpc == "" &&
					query.Address == ""
			},
			responses: []*pb.Stage{},
			output: [][]string{
				{
					NameText,
					PhaseText,
					AssetText,
					ServiceText,
					RpcText,
					AddressText,
				},
			},
		},
		{
			name: "no such service",
			args: []string{"--service", stageServiceForNum(4)},
			validateQuery: func(query *pb.Stage) bool {
				return query.Name == "" &&
					query.Phase == "" &&
					query.Asset == "" &&
					query.Service == stageServiceForNum(4) &&
					query.Rpc == "" &&
					query.Address == ""
			},
			responses: []*pb.Stage{},
			output: [][]string{
				{
					NameText,
					PhaseText,
					AssetText,
					ServiceText,
					RpcText,
					AddressText,
				},
			},
		},
		{
			name: "no such rpc",
			args: []string{"--rpc", stageRpcForNum(5)},
			validateQuery: func(query *pb.Stage) bool {
				return query.Name == "" &&
					query.Phase == "" &&
					query.Asset == "" &&
					query.Service == "" &&
					query.Rpc == stageRpcForNum(5) &&
					query.Address == ""
			},
			responses: []*pb.Stage{},
			output: [][]string{
				{
					NameText,
					PhaseText,
					AssetText,
					ServiceText,
					RpcText,
					AddressText,
				},
			},
		},
		{
			name: "no such address",
			args: []string{"--address", stageAddressForNum(6)},
			validateQuery: func(query *pb.Stage) bool {
				return query.Name == "" &&
					query.Phase == "" &&
					query.Asset == "" &&
					query.Service == "" &&
					query.Rpc == "" &&
					query.Address == stageAddressForNum(6)
			},
			responses: []*pb.Stage{},
			output: [][]string{
				{
					NameText,
					PhaseText,
					AssetText,
					ServiceText,
					RpcText,
					AddressText,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				lis := testutil.ListenAvailablePort(t)

				addr := lis.Addr().String()
				test.args = append(test.args, "--maestro", addr)

				mockServer := mock.MaestroServer{
					StageManagementServer: &mock.StageManagementServer{
						GetStageFn: func(
							query *pb.Stage,
							stream pb.StageManagement_GetServer,
						) error {
							if !test.validateQuery(query) {
								return fmt.Errorf(
									"validation failed with query %v",
									query)
							}
							for _, s := range test.responses {
								if err := stream.Send(s); err != nil {
									return fmt.Errorf("send failed: %v", err)
								}
							}
							return nil
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
				cmd := NewCmdGetStage()
				cmd.SetOut(b)
				cmd.SetArgs(test.args)
				err := cmd.Execute()
				assert.NilError(t, err, "execute error")
				out, err := ioutil.ReadAll(b)
				assert.NilError(t, err, "read output error")

				expectedOut, err := pterm.DefaultTable.
					WithHasHeader().
					WithData(test.output).
					Srender()
				expectedOut += "\n"
				assert.NilError(t, err, "render error")
				assert.Equal(t, expectedOut, string(out), "output differs")
			})
	}
}

func pbStageForNum(num int, phase apitypes.StagePhase) *pb.Stage {
	return &pb.Stage{
		Name:    stageNameForNum(num),
		Phase:   string(phase),
		Asset:   assetNameForNum(num),
		Service: stageServiceForNum(num),
		Rpc:     stageRpcForNum(num),
		Address: stageAddressForNum(num),
	}
}
