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
					query.Asset == "" &&
					query.Service == "" &&
					query.Method == "" &&
					query.Address == ""
			},
			responses: []*pb.Stage{},
			output: [][]string{
				{NameText, AssetText, ServiceText, MethodText, AddressText},
			},
		},
		{
			name: "one stage",
			args: []string{},
			validateQuery: func(query *pb.Stage) bool {
				return query.Name == "" &&
					query.Asset == "" &&
					query.Service == "" &&
					query.Method == "" &&
					query.Address == ""
			},
			responses: []*pb.Stage{pbStageForNum(0)},
			output: [][]string{
				{NameText, AssetText, ServiceText, MethodText, AddressText},
				{
					stageNameForNum(0),
					assetNameForNum(0),
					stageServiceForNum(0),
					stageMethodForNum(0),
					stageAddressForNum(0),
				},
			},
		},
		{
			name: "multiple stages",
			args: []string{},
			validateQuery: func(query *pb.Stage) bool {
				return query.Name == "" &&
					query.Asset == "" &&
					query.Service == "" &&
					query.Method == "" &&
					query.Address == ""
			},
			responses: []*pb.Stage{
				pbStageForNum(0),
				pbStageForNum(2),
				pbStageForNum(1),
			},
			output: [][]string{
				{NameText, AssetText, ServiceText, MethodText, AddressText},
				{
					stageNameForNum(0),
					assetNameForNum(0),
					stageServiceForNum(0),
					stageMethodForNum(0),
					stageAddressForNum(0),
				},
				{
					stageNameForNum(1),
					assetNameForNum(1),
					stageServiceForNum(1),
					stageMethodForNum(1),
					stageAddressForNum(1),
				},
				{
					stageNameForNum(2),
					assetNameForNum(2),
					stageServiceForNum(2),
					stageMethodForNum(2),
					stageAddressForNum(2),
				},
			},
		},
		{
			name: "filter by name",
			args: []string{stageNameForNum(2)},
			validateQuery: func(query *pb.Stage) bool {
				return query.Name == stageNameForNum(2) &&
					query.Asset == "" &&
					query.Service == "" &&
					query.Method == "" &&
					query.Address == ""
			},
			responses: []*pb.Stage{
				pbStageForNum(2),
			},
			output: [][]string{
				{NameText, AssetText, ServiceText, MethodText, AddressText},
				{
					stageNameForNum(2),
					assetNameForNum(2),
					stageServiceForNum(2),
					stageMethodForNum(2),
					stageAddressForNum(2),
				},
			},
		},
		{
			name: "filter by asset",
			args: []string{"--asset", assetNameForNum(2)},
			validateQuery: func(query *pb.Stage) bool {
				return query.Name == "" &&
					query.Asset == assetNameForNum(2) &&
					query.Service == "" &&
					query.Method == "" &&
					query.Address == ""
			},
			responses: []*pb.Stage{
				pbStageForNum(2),
			},
			output: [][]string{
				{NameText, AssetText, ServiceText, MethodText, AddressText},
				{
					stageNameForNum(2),
					assetNameForNum(2),
					stageServiceForNum(2),
					stageMethodForNum(2),
					stageAddressForNum(2),
				},
			},
		},
		{
			name: "filter by service",
			args: []string{"--service", stageServiceForNum(0)},
			validateQuery: func(query *pb.Stage) bool {
				return query.Name == "" &&
					query.Asset == "" &&
					query.Service == stageServiceForNum(0) &&
					query.Method == "" &&
					query.Address == ""
			},
			responses: []*pb.Stage{
				pbStageForNum(0),
			},
			output: [][]string{
				{NameText, AssetText, ServiceText, MethodText, AddressText},
				{
					stageNameForNum(0),
					assetNameForNum(0),
					stageServiceForNum(0),
					stageMethodForNum(0),
					stageAddressForNum(0),
				},
			},
		},
		{
			name: "filter by method",
			args: []string{"--method", stageMethodForNum(1)},
			validateQuery: func(query *pb.Stage) bool {
				return query.Name == "" &&
					query.Asset == "" &&
					query.Service == "" &&
					query.Method == stageMethodForNum(1) &&
					query.Address == ""
			},
			responses: []*pb.Stage{
				pbStageForNum(1),
			},
			output: [][]string{
				{NameText, AssetText, ServiceText, MethodText, AddressText},
				{
					stageNameForNum(1),
					assetNameForNum(1),
					stageServiceForNum(1),
					stageMethodForNum(1),
					stageAddressForNum(1),
				},
			},
		},
		{
			name: "no such name",
			args: []string{stageNameForNum(3)},
			validateQuery: func(query *pb.Stage) bool {
				return query.Name == stageNameForNum(3) &&
					query.Asset == "" &&
					query.Service == "" &&
					query.Method == "" &&
					query.Address == ""
			},
			responses: []*pb.Stage{},
			output: [][]string{
				{NameText, AssetText, ServiceText, MethodText, AddressText},
			},
		},
		{
			name: "no such asset",
			args: []string{"--asset", assetNameForNum(3)},
			validateQuery: func(query *pb.Stage) bool {
				return query.Name == "" &&
					query.Asset == assetNameForNum(3) &&
					query.Service == "" &&
					query.Method == "" &&
					query.Address == ""
			},
			responses: []*pb.Stage{},
			output: [][]string{
				{NameText, AssetText, ServiceText, MethodText, AddressText},
			},
		},
		{
			name: "no such service",
			args: []string{"--service", stageServiceForNum(4)},
			validateQuery: func(query *pb.Stage) bool {
				return query.Name == "" &&
					query.Asset == "" &&
					query.Service == stageServiceForNum(4) &&
					query.Method == "" &&
					query.Address == ""
			},
			responses: []*pb.Stage{},
			output: [][]string{
				{NameText, AssetText, ServiceText, MethodText, AddressText},
			},
		},
		{
			name: "no such method",
			args: []string{"--method", stageMethodForNum(5)},
			validateQuery: func(query *pb.Stage) bool {
				return query.Name == "" &&
					query.Asset == "" &&
					query.Service == "" &&
					query.Method == stageMethodForNum(5) &&
					query.Address == ""
			},
			responses: []*pb.Stage{},
			output: [][]string{
				{NameText, AssetText, ServiceText, MethodText, AddressText},
			},
		},
		{
			name: "no such address",
			args: []string{"--address", stageAddressForNum(6)},
			validateQuery: func(query *pb.Stage) bool {
				return query.Name == "" &&
					query.Asset == "" &&
					query.Service == "" &&
					query.Method == "" &&
					query.Address == stageAddressForNum(6)
			},
			responses: []*pb.Stage{},
			output: [][]string{
				{NameText, AssetText, ServiceText, MethodText, AddressText},
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

func pbStageForNum(num int) *pb.Stage {
	return &pb.Stage{
		Name:    stageNameForNum(num),
		Asset:   assetNameForNum(num),
		Service: stageServiceForNum(num),
		Method:  stageMethodForNum(num),
		Address: stageAddressForNum(num),
	}
}
