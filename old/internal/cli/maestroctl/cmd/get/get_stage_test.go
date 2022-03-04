package get

import (
	"bytes"
	"fmt"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/api"
	ipb "github.com/DuarteMRAlves/maestro/internal/api/pb"
	"github.com/pterm/pterm"
	"gotest.tools/v3/assert"
	"io/ioutil"
	"net"
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
					"stage-0",
					string(api.StageRunning),
					"asset-0",
					"service-0",
					"rpc-0",
					"address-0",
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
					"stage-0",
					string(api.StagePending),
					"asset-0",
					"service-0",
					"rpc-0",
					"address-0",
				},
				{
					"stage-1",
					string(api.StageFailed),
					"asset-1",
					"service-1",
					"rpc-1",
					"address-1",
				},
				{
					"stage-2",
					string(api.StageRunning),
					"asset-2",
					"service-2",
					"rpc-2",
					"address-2",
				},
			},
		},
		{
			name: "filter by name",
			args: []string{"stage-2"},
			validateReq: func(req *pb.GetStageRequest) bool {
				return req.Name == "stage-2" &&
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
					"stage-2",
					string(api.StageSucceeded),
					"asset-2",
					"service-2",
					"rpc-2",
					"address-2",
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
					"stage-1",
					string(api.StageRunning),
					"asset-1",
					"service-1",
					"rpc-1",
					"address-1",
				},
			},
		},
		{
			name: "filter by asset",
			args: []string{"--asset", "asset-2"},
			validateReq: func(req *pb.GetStageRequest) bool {
				return req.Name == "" &&
					req.Phase == "" &&
					req.Asset == "asset-2" &&
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
					"stage-2",
					string(api.StagePending),
					"asset-2",
					"service-2",
					"rpc-2",
					"address-2",
				},
			},
		},
		{
			name: "filter by service",
			args: []string{"--service", "service-0"},
			validateReq: func(req *pb.GetStageRequest) bool {
				return req.Name == "" &&
					req.Phase == "" &&
					req.Asset == "" &&
					req.Service == "service-0" &&
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
					"stage-0",
					string(api.StageRunning),
					"asset-0",
					"service-0",
					"rpc-0",
					"address-0",
				},
			},
		},
		{
			name: "filter by rpc",
			args: []string{"--rpc", "rpc-1"},
			validateReq: func(req *pb.GetStageRequest) bool {
				return req.Name == "" &&
					req.Phase == "" &&
					req.Asset == "" &&
					req.Service == "" &&
					req.Rpc == "rpc-1" &&
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
					"stage-1",
					string(api.StagePending),
					"asset-1",
					"service-1",
					"rpc-1",
					"address-1",
				},
			},
		},
		{
			name: "no such name",
			args: []string{"stage-3"},
			validateReq: func(req *pb.GetStageRequest) bool {
				return req.Name == "stage-3" &&
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
			args: []string{"--asset", "asset-3"},
			validateReq: func(req *pb.GetStageRequest) bool {
				return req.Name == "" &&
					req.Phase == "" &&
					req.Asset == "asset-3" &&
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
			args: []string{"--service", "service-4"},
			validateReq: func(req *pb.GetStageRequest) bool {
				return req.Name == "" &&
					req.Phase == "" &&
					req.Asset == "" &&
					req.Service == "service-4" &&
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
			args: []string{"--rpc", "rpc-5"},
			validateReq: func(req *pb.GetStageRequest) bool {
				return req.Name == "" &&
					req.Phase == "" &&
					req.Asset == "" &&
					req.Service == "" &&
					req.Rpc == "rpc-5" &&
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
			args: []string{"--address", "address-6"},
			validateReq: func(req *pb.GetStageRequest) bool {
				return req.Name == "" &&
					req.Phase == "" &&
					req.Asset == "" &&
					req.Service == "" &&
					req.Rpc == "" &&
					req.Address == "address-6"
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
				lis, err := net.Listen("tcp", "localhost:0")
				assert.NilError(t, err, "failed to listen")

				addr := lis.Addr().String()
				test.args = append(test.args, "--maestro", addr)

				mockServer := ipb.MockMaestroServer{
					ArchitectureManagementServer: &ipb.MockArchitectureManagementServer{
						GetStageFn: func(
							req *pb.GetStageRequest,
							stream pb.ArchitectureManagement_GetStageServer,
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
				err = cmd.Execute()
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
		Name:    fmt.Sprintf("stage-%d", num),
		Phase:   string(phase),
		Asset:   fmt.Sprintf("asset-%d", num),
		Service: fmt.Sprintf("service-%d", num),
		Rpc:     fmt.Sprintf("rpc-%d", num),
		Address: fmt.Sprintf("address-%d", num),
	}
}
