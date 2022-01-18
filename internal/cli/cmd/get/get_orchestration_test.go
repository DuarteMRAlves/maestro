package get

import (
	"bytes"
	"fmt"
	"github.com/DuarteMRAlves/maestro/api/pb"
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"github.com/DuarteMRAlves/maestro/internal/testutil"
	mockpb "github.com/DuarteMRAlves/maestro/internal/testutil/mock/pb"
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
		name          string
		args          []string
		validateQuery func(query *pb.Orchestration) bool
		responses     []*pb.Orchestration
		output        [][]string
	}{
		{
			name: "empty orchestrations",
			args: []string{},
			validateQuery: func(query *pb.Orchestration) bool {
				return query.Name == "" &&
					query.Phase == "" &&
					len(query.Links) == 0
			},
			responses: []*pb.Orchestration{},
			output:    [][]string{{NameText, PhaseText}},
		},
		{
			name: "one orchestration",
			args: []string{},
			validateQuery: func(query *pb.Orchestration) bool {
				return query.Name == "" &&
					query.Phase == "" &&
					len(query.Links) == 0
			},
			responses: []*pb.Orchestration{
				newPbOrchestration(0, apitypes.OrchestrationSucceeded),
			},
			output: [][]string{
				{NameText, PhaseText},
				{
					testutil.OrchestrationNameForNum(0),
					string(apitypes.OrchestrationSucceeded),
				},
			},
		},
		{
			name: "multiple orchestrations",
			args: []string{},
			validateQuery: func(query *pb.Orchestration) bool {
				return query.Name == "" &&
					query.Phase == "" &&
					len(query.Links) == 0
			},
			responses: []*pb.Orchestration{
				newPbOrchestration(1, apitypes.OrchestrationRunning),
				newPbOrchestration(0, apitypes.OrchestrationFailed),
				newPbOrchestration(2, apitypes.OrchestrationPending),
			},
			output: [][]string{
				{NameText, PhaseText},
				{
					testutil.OrchestrationNameForNum(0),
					string(apitypes.OrchestrationFailed),
				},
				{
					testutil.OrchestrationNameForNum(1),
					string(apitypes.OrchestrationRunning),
				},
				{
					testutil.OrchestrationNameForNum(2),
					string(apitypes.OrchestrationPending),
				},
			},
		},
		{
			name: "filter by name",
			args: []string{testutil.OrchestrationNameForNum(2)},
			validateQuery: func(query *pb.Orchestration) bool {
				return query.Name == testutil.OrchestrationNameForNum(2) &&
					query.Phase == "" &&
					len(query.Links) == 0
			},
			responses: []*pb.Orchestration{
				newPbOrchestration(2, apitypes.OrchestrationSucceeded),
			},
			output: [][]string{
				{NameText, PhaseText},
				{
					testutil.OrchestrationNameForNum(2),
					string(apitypes.OrchestrationSucceeded),
				},
			},
		},
		{
			name: "no such name",
			args: []string{testutil.OrchestrationNameForNum(3)},
			validateQuery: func(query *pb.Orchestration) bool {
				return query.Name == testutil.OrchestrationNameForNum(3) &&
					query.Phase == "" &&
					len(query.Links) == 0
			},
			responses: []*pb.Orchestration{},
			output:    [][]string{{NameText, PhaseText}},
		},
		{
			name: "filter by phase",
			args: []string{"--phase", string(apitypes.OrchestrationPending)},
			validateQuery: func(query *pb.Orchestration) bool {
				return query.Name == "" &&
					query.Phase == string(apitypes.OrchestrationPending) &&
					len(query.Links) == 0
			},
			responses: []*pb.Orchestration{
				newPbOrchestration(1, apitypes.OrchestrationPending),
			},
			output: [][]string{
				{NameText, PhaseText},
				{
					testutil.OrchestrationNameForNum(1),
					string(apitypes.OrchestrationPending),
				},
			},
		},
		{
			name: "no such phase",
			args: []string{"--phase", string(apitypes.OrchestrationRunning)},
			validateQuery: func(query *pb.Orchestration) bool {
				return query.Name == "" &&
					query.Phase == string(apitypes.OrchestrationRunning) &&
					len(query.Links) == 0
			},
			responses: []*pb.Orchestration{},
			output:    [][]string{{NameText, PhaseText}},
		},
	}
	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				lis := testutil.ListenAvailablePort(t)

				addr := lis.Addr().String()
				test.args = append(test.args, "--maestro", addr)

				mockServer := mockpb.MaestroServer{
					OrchestrationManagementServer: &mockpb.OrchestrationManagementServer{
						GetOrchestrationFn: func(
							query *pb.Orchestration,
							stream pb.OrchestrationManagement_GetServer,
						) error {
							if !test.validateQuery(query) {
								return fmt.Errorf(
									"validation failed with query %v",
									query)
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
			})
	}
}

func newPbOrchestration(
	num int,
	phase apitypes.OrchestrationPhase,
) *pb.Orchestration {
	return &pb.Orchestration{
		Name:  testutil.OrchestrationNameForNum(num),
		Phase: string(phase),
		Links: nil,
	}
}
