package get

import (
	"bytes"
	"fmt"
	"github.com/DuarteMRAlves/maestro/api/pb"
	ipb "github.com/DuarteMRAlves/maestro/internal/api/pb"
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"github.com/DuarteMRAlves/maestro/internal/testutil"
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
					testutil.StageNameForNumStr(0),
					string(apitypes.StageRunning),
					testutil.AssetNameForNumStr(0),
					testutil.StageServiceForNum(0),
					testutil.StageRpcForNum(0),
					testutil.StageAddressForNum(0),
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
					testutil.StageNameForNumStr(0),
					string(apitypes.StagePending),
					testutil.AssetNameForNumStr(0),
					testutil.StageServiceForNum(0),
					testutil.StageRpcForNum(0),
					testutil.StageAddressForNum(0),
				},
				{
					testutil.StageNameForNumStr(1),
					string(apitypes.StageFailed),
					testutil.AssetNameForNumStr(1),
					testutil.StageServiceForNum(1),
					testutil.StageRpcForNum(1),
					testutil.StageAddressForNum(1),
				},
				{
					testutil.StageNameForNumStr(2),
					string(apitypes.StageRunning),
					testutil.AssetNameForNumStr(2),
					testutil.StageServiceForNum(2),
					testutil.StageRpcForNum(2),
					testutil.StageAddressForNum(2),
				},
			},
		},
		{
			name: "filter by name",
			args: []string{testutil.StageNameForNumStr(2)},
			validateQuery: func(query *pb.Stage) bool {
				return query.Name == testutil.StageNameForNumStr(2) &&
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
					testutil.StageNameForNumStr(2),
					string(apitypes.StageSucceeded),
					testutil.AssetNameForNumStr(2),
					testutil.StageServiceForNum(2),
					testutil.StageRpcForNum(2),
					testutil.StageAddressForNum(2),
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
					testutil.StageNameForNumStr(1),
					string(apitypes.StageRunning),
					testutil.AssetNameForNumStr(1),
					testutil.StageServiceForNum(1),
					testutil.StageRpcForNum(1),
					testutil.StageAddressForNum(1),
				},
			},
		},
		{
			name: "filter by asset",
			args: []string{"--asset", testutil.AssetNameForNumStr(2)},
			validateQuery: func(query *pb.Stage) bool {
				return query.Name == "" &&
					query.Phase == "" &&
					query.Asset == testutil.AssetNameForNumStr(2) &&
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
					testutil.StageNameForNumStr(2),
					string(apitypes.StagePending),
					testutil.AssetNameForNumStr(2),
					testutil.StageServiceForNum(2),
					testutil.StageRpcForNum(2),
					testutil.StageAddressForNum(2),
				},
			},
		},
		{
			name: "filter by service",
			args: []string{"--service", testutil.StageServiceForNum(0)},
			validateQuery: func(query *pb.Stage) bool {
				return query.Name == "" &&
					query.Phase == "" &&
					query.Asset == "" &&
					query.Service == testutil.StageServiceForNum(0) &&
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
					testutil.StageNameForNumStr(0),
					string(apitypes.StageRunning),
					testutil.AssetNameForNumStr(0),
					testutil.StageServiceForNum(0),
					testutil.StageRpcForNum(0),
					testutil.StageAddressForNum(0),
				},
			},
		},
		{
			name: "filter by rpc",
			args: []string{"--rpc", testutil.StageRpcForNum(1)},
			validateQuery: func(query *pb.Stage) bool {
				return query.Name == "" &&
					query.Phase == "" &&
					query.Asset == "" &&
					query.Service == "" &&
					query.Rpc == testutil.StageRpcForNum(1) &&
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
					testutil.StageNameForNumStr(1),
					string(apitypes.StagePending),
					testutil.AssetNameForNumStr(1),
					testutil.StageServiceForNum(1),
					testutil.StageRpcForNum(1),
					testutil.StageAddressForNum(1),
				},
			},
		},
		{
			name: "no such name",
			args: []string{testutil.StageNameForNumStr(3)},
			validateQuery: func(query *pb.Stage) bool {
				return query.Name == testutil.StageNameForNumStr(3) &&
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
			args: []string{"--asset", testutil.AssetNameForNumStr(3)},
			validateQuery: func(query *pb.Stage) bool {
				return query.Name == "" &&
					query.Phase == "" &&
					query.Asset == testutil.AssetNameForNumStr(3) &&
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
			args: []string{"--service", testutil.StageServiceForNum(4)},
			validateQuery: func(query *pb.Stage) bool {
				return query.Name == "" &&
					query.Phase == "" &&
					query.Asset == "" &&
					query.Service == testutil.StageServiceForNum(4) &&
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
			args: []string{"--rpc", testutil.StageRpcForNum(5)},
			validateQuery: func(query *pb.Stage) bool {
				return query.Name == "" &&
					query.Phase == "" &&
					query.Asset == "" &&
					query.Service == "" &&
					query.Rpc == testutil.StageRpcForNum(5) &&
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
			args: []string{"--address", testutil.StageAddressForNum(6)},
			validateQuery: func(query *pb.Stage) bool {
				return query.Name == "" &&
					query.Phase == "" &&
					query.Asset == "" &&
					query.Service == "" &&
					query.Rpc == "" &&
					query.Address == testutil.StageAddressForNum(6)
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

				mockServer := ipb.MockMaestroServer{
					StageManagementServer: &ipb.MockStageManagementServer{
						GetStageFn: func(
							query *pb.Stage,
							stream pb.StageManagement_GetServer,
						) error {
							if !test.validateQuery(query) {
								return fmt.Errorf(
									"validation failed with query %v",
									query,
								)
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
			},
		)
	}
}

func pbStageForNum(num int, phase apitypes.StagePhase) *pb.Stage {
	return &pb.Stage{
		Name:    testutil.StageNameForNumStr(num),
		Phase:   string(phase),
		Asset:   testutil.AssetNameForNumStr(num),
		Service: testutil.StageServiceForNum(num),
		Rpc:     testutil.StageRpcForNum(num),
		Address: testutil.StageAddressForNum(num),
	}
}
