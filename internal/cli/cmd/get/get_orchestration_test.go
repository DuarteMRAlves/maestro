package get

import (
	"bytes"
	"fmt"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/testutil"
	"github.com/DuarteMRAlves/maestro/internal/testutil/mock"
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
				return query.Name == "" && len(query.Links) == 0
			},
			responses: []*pb.Orchestration{},
			output:    [][]string{{NameText}},
		},
		{
			name: "one orchestration",
			args: []string{},
			validateQuery: func(query *pb.Orchestration) bool {
				return query.Name == "" && len(query.Links) == 0
			},
			responses: []*pb.Orchestration{pbOrchestrationForNum(0)},
			output:    [][]string{{NameText}, {orchestrationNameForNum(0)}},
		},
		{
			name: "multiple orchestrations",
			args: []string{},
			validateQuery: func(query *pb.Orchestration) bool {
				return query.Name == "" && len(query.Links) == 0
			},
			responses: []*pb.Orchestration{
				pbOrchestrationForNum(1),
				pbOrchestrationForNum(0),
				pbOrchestrationForNum(2),
			},
			output: [][]string{
				{NameText},
				{orchestrationNameForNum(0)},
				{orchestrationNameForNum(1)},
				{orchestrationNameForNum(2)},
			},
		},
		{
			name: "filter by name",
			args: []string{orchestrationNameForNum(2)},
			validateQuery: func(query *pb.Orchestration) bool {
				return query.Name == orchestrationNameForNum(2) &&
					len(query.Links) == 0
			},
			responses: []*pb.Orchestration{pbOrchestrationForNum(2)},
			output:    [][]string{{NameText}, {orchestrationNameForNum(2)}},
		},
		{
			name: "no such name",
			args: []string{orchestrationNameForNum(3)},
			validateQuery: func(query *pb.Orchestration) bool {
				return query.Name == orchestrationNameForNum(3) &&
					len(query.Links) == 0
			},
			responses: []*pb.Orchestration{},
			output:    [][]string{{NameText}},
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
					OrchestrationManagementServer: &mock.OrchestrationManagementServer{
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

func pbOrchestrationForNum(num int) *pb.Orchestration {
	return &pb.Orchestration{Name: orchestrationNameForNum(num), Links: nil}
}
