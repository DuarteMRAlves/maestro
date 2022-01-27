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
				{
					util.OrchestrationNameForNumStr(0),
					string(api.OrchestrationSucceeded),
				},
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
				{
					util.OrchestrationNameForNumStr(0),
					string(api.OrchestrationFailed),
				},
				{
					util.OrchestrationNameForNumStr(1),
					string(api.OrchestrationRunning),
				},
				{
					util.OrchestrationNameForNumStr(2),
					string(api.OrchestrationPending),
				},
			},
		},
		{
			name: "filter by name",
			args: []string{util.OrchestrationNameForNumStr(2)},
			validateReq: func(req *pb.GetOrchestrationRequest) bool {
				return req.Name == util.OrchestrationNameForNumStr(2) &&
					req.Phase == ""
			},
			responses: []*pb.Orchestration{
				newPbOrchestration(2, api.OrchestrationSucceeded),
			},
			output: [][]string{
				{NameText, PhaseText},
				{
					util.OrchestrationNameForNumStr(2),
					string(api.OrchestrationSucceeded),
				},
			},
		},
		{
			name: "no such name",
			args: []string{util.OrchestrationNameForNumStr(3)},
			validateReq: func(req *pb.GetOrchestrationRequest) bool {
				return req.Name == util.OrchestrationNameForNumStr(3) &&
					req.Phase == ""
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
				{
					util.OrchestrationNameForNumStr(1),
					string(api.OrchestrationPending),
				},
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
				lis := util.NewTestListener(t)

				addr := lis.Addr().String()
				test.args = append(test.args, "--maestro", addr)

				mockServer := ipb.MockMaestroServer{
					OrchestrationManagementServer: &ipb.MockOrchestrationManagementServer{
						GetOrchestrationFn: func(
							req *pb.GetOrchestrationRequest,
							stream pb.OrchestrationManagement_GetServer,
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

func newPbOrchestration(
	num int,
	phase api.OrchestrationPhase,
) *pb.Orchestration {
	return &pb.Orchestration{
		Name:  util.OrchestrationNameForNumStr(num),
		Phase: string(phase),
		Links: nil,
	}
}
