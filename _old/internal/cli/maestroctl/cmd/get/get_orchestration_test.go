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

// TestGetOrchestration_CorrectDisplay performs integration testing on the
// GetOrchestration command considering operations that produce table outputs.
// It runs a maestro server and then executes a get orchestration command with
// predetermined arguments, verifying its output by comparing with an expected
// table.
func TestGetOrchestration_CorrectDisplay(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		validateReq func(req *pb.GetOrchestrationRequest) bool
		responses   []*pb.Orchestration
		output      [][]string
	}{
		{
			name: "empty orchestrations",
			args: []string{},
			validateReq: func(req *pb.GetOrchestrationRequest) bool {
				return req.Name == "" &&
					req.Phase == ""
			},
			responses: []*pb.Orchestration{},
			output:    [][]string{{NameText, PhaseText}},
		},
		{
			name: "one orchestration",
			args: []string{},
			validateReq: func(req *pb.GetOrchestrationRequest) bool {
				return req.Name == "" &&
					req.Phase == ""
			},
			responses: []*pb.Orchestration{
				newPbOrchestration(0, api.OrchestrationSucceeded),
			},
			output: [][]string{
				{NameText, PhaseText},
				{"orchestration-0", string(api.OrchestrationSucceeded)},
			},
		},
		{
			name: "multiple orchestrations",
			args: []string{},
			validateReq: func(req *pb.GetOrchestrationRequest) bool {
				return req.Name == "" &&
					req.Phase == ""
			},
			responses: []*pb.Orchestration{
				newPbOrchestration(1, api.OrchestrationRunning),
				newPbOrchestration(0, api.OrchestrationFailed),
				newPbOrchestration(2, api.OrchestrationPending),
			},
			output: [][]string{
				{NameText, PhaseText},
				{"orchestration-0", string(api.OrchestrationFailed)},
				{"orchestration-1", string(api.OrchestrationRunning)},
				{"orchestration-2", string(api.OrchestrationPending)},
			},
		},
		{
			name: "filter by name",
			args: []string{"orchestration-2"},
			validateReq: func(req *pb.GetOrchestrationRequest) bool {
				return req.Name == "orchestration-2" &&
					req.Phase == ""
			},
			responses: []*pb.Orchestration{
				newPbOrchestration(2, api.OrchestrationSucceeded),
			},
			output: [][]string{
				{NameText, PhaseText},
				{
					"orchestration-2",
					string(api.OrchestrationSucceeded),
				},
			},
		},
		{
			name: "no such name",
			args: []string{"orchestration-3"},
			validateReq: func(req *pb.GetOrchestrationRequest) bool {
				return req.Name == "orchestration-3" && req.Phase == ""
			},
			responses: []*pb.Orchestration{},
			output:    [][]string{{NameText, PhaseText}},
		},
		{
			name: "filter by phase",
			args: []string{"--phase", string(api.OrchestrationPending)},
			validateReq: func(req *pb.GetOrchestrationRequest) bool {
				return req.Name == "" &&
					req.Phase == string(api.OrchestrationPending)
			},
			responses: []*pb.Orchestration{
				newPbOrchestration(1, api.OrchestrationPending),
			},
			output: [][]string{
				{NameText, PhaseText},
				{"orchestration-1", string(api.OrchestrationPending)},
			},
		},
		{
			name: "no such phase",
			args: []string{"--phase", string(api.OrchestrationRunning)},
			validateReq: func(req *pb.GetOrchestrationRequest) bool {
				return req.Name == "" &&
					req.Phase == string(api.OrchestrationRunning)
			},
			responses: []*pb.Orchestration{},
			output:    [][]string{{NameText, PhaseText}},
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
						GetOrchestrationFn: func(
							req *pb.GetOrchestrationRequest,
							stream pb.ArchitectureManagement_GetOrchestrationServer,
						) error {
							if !test.validateReq(req) {
								return fmt.Errorf(
									"validation failed with req %v",
									req,
								)
							}
							for _, o := range test.responses {
								if err := stream.Send(o); err != nil {
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
				cmd := NewCmdGetOrchestration()
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

func newPbOrchestration(
	num int,
	phase api.OrchestrationPhase,
) *pb.Orchestration {
	return &pb.Orchestration{
		Name:  fmt.Sprintf("orchestration-%d", num),
		Phase: string(phase),
		Links: nil,
	}
}