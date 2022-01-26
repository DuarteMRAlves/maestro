package get

import (
	"bytes"
	"fmt"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/api"
	ipb "github.com/DuarteMRAlves/maestro/internal/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/util"
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
		name        string
		args        []string
		validateReq func(*pb.GetStageRequest) bool
		responses   []*pb.Stage
		output      [][]string
	}{
		{
			name: "empty stages",
			args: []string{},
			validateReq: func(req *pb.GetStageRequest) bool {
				return req.Name == "" &&
					req.Phase == "" &&
					req.Asset == "" &&
					req.Service == "" &&
					req.Rpc == "" &&
					req.Address == ""
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
			validateReq: func(req *pb.GetStageRequest) bool {
				return req.Name == "" &&
					req.Phase == "" &&
					req.Asset == "" &&
					req.Service == "" &&
					req.Rpc == "" &&
					req.Address == ""
			},
			responses: []*pb.Stage{pbStageForNum(0, api.StageRunning)},
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
					util.StageNameForNumStr(0),
					string(api.StageRunning),
					util.AssetNameForNumStr(0),
					util.StageServiceForNum(0),
					util.StageRpcForNum(0),
					util.StageAddressForNum(0),
				},
			},
		},
		{
			name: "multiple stages",
			args: []string{},
			validateReq: func(req *pb.GetStageRequest) bool {
				return req.Name == "" &&
					req.Phase == "" &&
					req.Asset == "" &&
					req.Service == "" &&
					req.Rpc == "" &&
					req.Address == ""
			},
			responses: []*pb.Stage{
				pbStageForNum(0, api.StagePending),
				pbStageForNum(2, api.StageRunning),
				pbStageForNum(1, api.StageFailed),
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
					util.StageNameForNumStr(0),
					string(api.StagePending),
					util.AssetNameForNumStr(0),
					util.StageServiceForNum(0),
					util.StageRpcForNum(0),
					util.StageAddressForNum(0),
				},
				{
					util.StageNameForNumStr(1),
					string(api.StageFailed),
					util.AssetNameForNumStr(1),
					util.StageServiceForNum(1),
					util.StageRpcForNum(1),
					util.StageAddressForNum(1),
				},
				{
					util.StageNameForNumStr(2),
					string(api.StageRunning),
					util.AssetNameForNumStr(2),
					util.StageServiceForNum(2),
					util.StageRpcForNum(2),
					util.StageAddressForNum(2),
				},
			},
		},
		{
			name: "filter by name",
			args: []string{util.StageNameForNumStr(2)},
			validateReq: func(req *pb.GetStageRequest) bool {
				return req.Name == util.StageNameForNumStr(2) &&
					req.Phase == "" &&
					req.Asset == "" &&
					req.Service == "" &&
					req.Rpc == "" &&
					req.Address == ""
			},
			responses: []*pb.Stage{
				pbStageForNum(2, api.StageSucceeded),
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
					util.StageNameForNumStr(2),
					string(api.StageSucceeded),
					util.AssetNameForNumStr(2),
					util.StageServiceForNum(2),
					util.StageRpcForNum(2),
					util.StageAddressForNum(2),
				},
			},
		},
		{
			name: "filter by phase",
			args: []string{"--phase", string(api.StageRunning)},
			validateReq: func(req *pb.GetStageRequest) bool {
				return req.Name == "" &&
					req.Phase == string(api.StageRunning) &&
					req.Asset == "" &&
					req.Service == "" &&
					req.Rpc == "" &&
					req.Address == ""
			},
			responses: []*pb.Stage{
				pbStageForNum(1, api.StageRunning),
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
					util.StageNameForNumStr(1),
					string(api.StageRunning),
					util.AssetNameForNumStr(1),
					util.StageServiceForNum(1),
					util.StageRpcForNum(1),
					util.StageAddressForNum(1),
				},
			},
		},
		{
			name: "filter by asset",
			args: []string{"--asset", util.AssetNameForNumStr(2)},
			validateReq: func(req *pb.GetStageRequest) bool {
				return req.Name == "" &&
					req.Phase == "" &&
					req.Asset == util.AssetNameForNumStr(2) &&
					req.Service == "" &&
					req.Rpc == "" &&
					req.Address == ""
			},
			responses: []*pb.Stage{
				pbStageForNum(2, api.StagePending),
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
					util.StageNameForNumStr(2),
					string(api.StagePending),
					util.AssetNameForNumStr(2),
					util.StageServiceForNum(2),
					util.StageRpcForNum(2),
					util.StageAddressForNum(2),
				},
			},
		},
		{
			name: "filter by service",
			args: []string{"--service", util.StageServiceForNum(0)},
			validateReq: func(req *pb.GetStageRequest) bool {
				return req.Name == "" &&
					req.Phase == "" &&
					req.Asset == "" &&
					req.Service == util.StageServiceForNum(0) &&
					req.Rpc == "" &&
					req.Address == ""
			},
			responses: []*pb.Stage{
				pbStageForNum(0, api.StageRunning),
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
					util.StageNameForNumStr(0),
					string(api.StageRunning),
					util.AssetNameForNumStr(0),
					util.StageServiceForNum(0),
					util.StageRpcForNum(0),
					util.StageAddressForNum(0),
				},
			},
		},
		{
			name: "filter by rpc",
			args: []string{"--rpc", util.StageRpcForNum(1)},
			validateReq: func(req *pb.GetStageRequest) bool {
				return req.Name == "" &&
					req.Phase == "" &&
					req.Asset == "" &&
					req.Service == "" &&
					req.Rpc == util.StageRpcForNum(1) &&
					req.Address == ""
			},
			responses: []*pb.Stage{
				pbStageForNum(1, api.StagePending),
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
					util.StageNameForNumStr(1),
					string(api.StagePending),
					util.AssetNameForNumStr(1),
					util.StageServiceForNum(1),
					util.StageRpcForNum(1),
					util.StageAddressForNum(1),
				},
			},
		},
		{
			name: "no such name",
			args: []string{util.StageNameForNumStr(3)},
			validateReq: func(req *pb.GetStageRequest) bool {
				return req.Name == util.StageNameForNumStr(3) &&
					req.Phase == "" &&
					req.Asset == "" &&
					req.Service == "" &&
					req.Rpc == "" &&
					req.Address == ""
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
			args: []string{"--phase", string(api.StagePending)},
			validateReq: func(req *pb.GetStageRequest) bool {
				return req.Name == "" &&
					req.Phase == string(api.StagePending) &&
					req.Asset == "" &&
					req.Service == "" &&
					req.Rpc == "" &&
					req.Address == ""
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
			args: []string{"--asset", util.AssetNameForNumStr(3)},
			validateReq: func(req *pb.GetStageRequest) bool {
				return req.Name == "" &&
					req.Phase == "" &&
					req.Asset == util.AssetNameForNumStr(3) &&
					req.Service == "" &&
					req.Rpc == "" &&
					req.Address == ""
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
			args: []string{"--service", util.StageServiceForNum(4)},
			validateReq: func(req *pb.GetStageRequest) bool {
				return req.Name == "" &&
					req.Phase == "" &&
					req.Asset == "" &&
					req.Service == util.StageServiceForNum(4) &&
					req.Rpc == "" &&
					req.Address == ""
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
			args: []string{"--rpc", util.StageRpcForNum(5)},
			validateReq: func(req *pb.GetStageRequest) bool {
				return req.Name == "" &&
					req.Phase == "" &&
					req.Asset == "" &&
					req.Service == "" &&
					req.Rpc == util.StageRpcForNum(5) &&
					req.Address == ""
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
			args: []string{"--address", util.StageAddressForNum(6)},
			validateReq: func(req *pb.GetStageRequest) bool {
				return req.Name == "" &&
					req.Phase == "" &&
					req.Asset == "" &&
					req.Service == "" &&
					req.Rpc == "" &&
					req.Address == util.StageAddressForNum(6)
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
				lis := util.NewTestListener(t)

				addr := lis.Addr().String()
				test.args = append(test.args, "--maestro", addr)

				mockServer := ipb.MockMaestroServer{
					StageManagementServer: &ipb.MockStageManagementServer{
						GetStageFn: func(
							req *pb.GetStageRequest,
							stream pb.StageManagement_GetServer,
						) error {
							if !test.validateReq(req) {
								return fmt.Errorf(
									"validation failed with req %v",
									req,
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

func pbStageForNum(num int, phase api.StagePhase) *pb.Stage {
	return &pb.Stage{
		Name:    util.StageNameForNumStr(num),
		Phase:   string(phase),
		Asset:   util.AssetNameForNumStr(num),
		Service: util.StageServiceForNum(num),
		Rpc:     util.StageRpcForNum(num),
		Address: util.StageAddressForNum(num),
	}
}
