package get

import (
	"bytes"
	"fmt"
	"github.com/DuarteMRAlves/maestro/api/pb"
	ipb "github.com/DuarteMRAlves/maestro/internal/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/testutil"
	"github.com/pterm/pterm"
	"gotest.tools/v3/assert"
	"io/ioutil"
	"testing"
)

// TestGetLink_CorrectDisplay performs testing on the GetLink command
// considering operations that produce table outputs. It runs a mock maestro
// server and then executes a get link command with predetermined arguments,
// verifying its output by comparing with an expected table.
func TestGetLink_CorrectDisplay(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		validateReq func(*pb.GetLinkRequest) bool
		responses   []*pb.Link
		output      [][]string
	}{
		{
			name: "empty links",
			args: []string{},
			validateReq: func(req *pb.GetLinkRequest) bool {
				return req.Name == "" &&
					req.SourceStage == "" &&
					req.SourceField == "" &&
					req.TargetStage == "" &&
					req.TargetField == ""
			},
			responses: []*pb.Link{},
			output: [][]string{
				{
					NameText,
					SourceStageText,
					SourceFieldText,
					TargetStageText,
					TargetFieldText,
				},
			},
		},
		{
			name: "one link",
			args: []string{},
			validateReq: func(req *pb.GetLinkRequest) bool {
				return req.Name == "" &&
					req.SourceStage == "" &&
					req.SourceField == "" &&
					req.TargetStage == "" &&
					req.TargetField == ""
			},
			responses: []*pb.Link{pbLinkForNum(0)},
			output: [][]string{
				{
					NameText,
					SourceStageText,
					SourceFieldText,
					TargetStageText,
					TargetFieldText,
				},
				{
					testutil.LinkNameForNumStr(0),
					testutil.LinkSourceStageForNumStr(0),
					testutil.LinkSourceFieldForNum(0),
					testutil.LinkTargetStageForNumStr(0),
					testutil.LinkTargetFieldForNum(0),
				},
			},
		},
		{
			name: "multiple links",
			args: []string{},
			validateReq: func(req *pb.GetLinkRequest) bool {
				return req.Name == "" &&
					req.SourceStage == "" &&
					req.SourceField == "" &&
					req.TargetStage == "" &&
					req.TargetField == ""
			},
			responses: []*pb.Link{
				pbLinkForNum(1),
				pbLinkForNum(2),
				pbLinkForNum(0),
			},
			output: [][]string{
				{
					NameText,
					SourceStageText,
					SourceFieldText,
					TargetStageText,
					TargetFieldText,
				},
				{
					testutil.LinkNameForNumStr(0),
					testutil.LinkSourceStageForNumStr(0),
					testutil.LinkSourceFieldForNum(0),
					testutil.LinkTargetStageForNumStr(0),
					testutil.LinkTargetFieldForNum(0),
				},
				{
					testutil.LinkNameForNumStr(1),
					testutil.LinkSourceStageForNumStr(1),
					testutil.LinkSourceFieldForNum(1),
					testutil.LinkTargetStageForNumStr(1),
					testutil.LinkTargetFieldForNum(1),
				},
				{
					testutil.LinkNameForNumStr(2),
					testutil.LinkSourceStageForNumStr(2),
					testutil.LinkSourceFieldForNum(2),
					testutil.LinkTargetStageForNumStr(2),
					testutil.LinkTargetFieldForNum(2),
				},
			},
		},
		{
			name: "filter by name",
			args: []string{testutil.LinkNameForNumStr(2)},
			validateReq: func(req *pb.GetLinkRequest) bool {
				return req.Name == testutil.LinkNameForNumStr(2) &&
					req.SourceStage == "" &&
					req.SourceField == "" &&
					req.TargetStage == "" &&
					req.TargetField == ""
			},
			responses: []*pb.Link{pbLinkForNum(2)},
			output: [][]string{
				{
					NameText,
					SourceStageText,
					SourceFieldText,
					TargetStageText,
					TargetFieldText,
				},
				{
					testutil.LinkNameForNumStr(2),
					testutil.LinkSourceStageForNumStr(2),
					testutil.LinkSourceFieldForNum(2),
					testutil.LinkTargetStageForNumStr(2),
					testutil.LinkTargetFieldForNum(2),
				},
			},
		},
		{
			name: "filter by source stage",
			args: []string{
				"--source-stage",
				testutil.LinkSourceStageForNumStr(2),
			},
			validateReq: func(req *pb.GetLinkRequest) bool {
				return req.Name == "" &&
					req.SourceStage == testutil.LinkSourceStageForNumStr(2) &&
					req.SourceField == "" &&
					req.TargetStage == "" &&
					req.TargetField == ""
			},
			responses: []*pb.Link{pbLinkForNum(2)},
			output: [][]string{
				{
					NameText,
					SourceStageText,
					SourceFieldText,
					TargetStageText,
					TargetFieldText,
				},
				{
					testutil.LinkNameForNumStr(2),
					testutil.LinkSourceStageForNumStr(2),
					testutil.LinkSourceFieldForNum(2),
					testutil.LinkTargetStageForNumStr(2),
					testutil.LinkTargetFieldForNum(2),
				},
			},
		},
		{
			name: "filter by source field",
			args: []string{"--source-field", testutil.LinkSourceFieldForNum(0)},
			validateReq: func(req *pb.GetLinkRequest) bool {
				return req.Name == "" &&
					req.SourceStage == "" &&
					req.SourceField == testutil.LinkSourceFieldForNum(0) &&
					req.TargetStage == "" &&
					req.TargetField == ""
			},
			responses: []*pb.Link{pbLinkForNum(0)},
			output: [][]string{
				{
					NameText,
					SourceStageText,
					SourceFieldText,
					TargetStageText,
					TargetFieldText,
				},
				{
					testutil.LinkNameForNumStr(0),
					testutil.LinkSourceStageForNumStr(0),
					testutil.LinkSourceFieldForNum(0),
					testutil.LinkTargetStageForNumStr(0),
					testutil.LinkTargetFieldForNum(0),
				},
			},
		},
		{
			name: "filter by target stage",
			args: []string{
				"--target-stage",
				testutil.LinkTargetStageForNumStr(1),
			},
			validateReq: func(req *pb.GetLinkRequest) bool {
				return req.Name == "" &&
					req.SourceStage == "" &&
					req.SourceField == "" &&
					req.TargetStage == testutil.LinkTargetStageForNumStr(1) &&
					req.TargetField == ""
			},
			responses: []*pb.Link{pbLinkForNum(1)},
			output: [][]string{
				{
					NameText,
					SourceStageText,
					SourceFieldText,
					TargetStageText,
					TargetFieldText,
				},
				{
					testutil.LinkNameForNumStr(1),
					testutil.LinkSourceStageForNumStr(1),
					testutil.LinkSourceFieldForNum(1),
					testutil.LinkTargetStageForNumStr(1),
					testutil.LinkTargetFieldForNum(1),
				},
			},
		},
		{
			name: "filter by target field",
			args: []string{"--target-field", testutil.LinkTargetFieldForNum(2)},
			validateReq: func(req *pb.GetLinkRequest) bool {
				return req.Name == "" &&
					req.SourceStage == "" &&
					req.SourceField == "" &&
					req.TargetStage == "" &&
					req.TargetField == testutil.LinkTargetFieldForNum(2)
			},
			responses: []*pb.Link{pbLinkForNum(2)},
			output: [][]string{
				{
					NameText,
					SourceStageText,
					SourceFieldText,
					TargetStageText,
					TargetFieldText,
				},
				{
					testutil.LinkNameForNumStr(2),
					testutil.LinkSourceStageForNumStr(2),
					testutil.LinkSourceFieldForNum(2),
					testutil.LinkTargetStageForNumStr(2),
					testutil.LinkTargetFieldForNum(2),
				},
			},
		},
		{
			name: "no such name",
			args: []string{testutil.LinkNameForNumStr(3)},
			validateReq: func(req *pb.GetLinkRequest) bool {
				return req.Name == testutil.LinkNameForNumStr(3) &&
					req.SourceStage == "" &&
					req.SourceField == "" &&
					req.TargetStage == "" &&
					req.TargetField == ""
			},
			responses: []*pb.Link{},
			output: [][]string{
				{
					NameText,
					SourceStageText,
					SourceFieldText,
					TargetStageText,
					TargetFieldText,
				},
			},
		},
		{
			name: "no such source stage",
			args: []string{
				"--source-stage",
				testutil.LinkSourceStageForNumStr(3),
			},
			validateReq: func(req *pb.GetLinkRequest) bool {
				return req.Name == "" &&
					req.SourceStage == testutil.LinkSourceStageForNumStr(3) &&
					req.SourceField == "" &&
					req.TargetStage == "" &&
					req.TargetField == ""
			},
			responses: []*pb.Link{},
			output: [][]string{
				{
					NameText,
					SourceStageText,
					SourceFieldText,
					TargetStageText,
					TargetFieldText,
				},
			},
		},
		{
			name: "no such source field",
			args: []string{"--source-field", testutil.LinkSourceFieldForNum(4)},
			validateReq: func(req *pb.GetLinkRequest) bool {
				return req.Name == "" &&
					req.SourceStage == "" &&
					req.SourceField == testutil.LinkSourceFieldForNum(4) &&
					req.TargetStage == "" &&
					req.TargetField == ""
			},
			responses: []*pb.Link{},
			output: [][]string{
				{
					NameText,
					SourceStageText,
					SourceFieldText,
					TargetStageText,
					TargetFieldText,
				},
			},
		},
		{
			name: "no such target stage",
			args: []string{
				"--target-stage",
				testutil.LinkTargetStageForNumStr(5),
			},
			validateReq: func(req *pb.GetLinkRequest) bool {
				return req.Name == "" &&
					req.SourceStage == "" &&
					req.SourceField == "" &&
					req.TargetStage == testutil.LinkTargetStageForNumStr(5) &&
					req.TargetField == ""
			},
			responses: []*pb.Link{},
			output: [][]string{
				{
					NameText,
					SourceStageText,
					SourceFieldText,
					TargetStageText,
					TargetFieldText,
				},
			},
		},
		{
			name: "no such target field",
			args: []string{"--target-field", testutil.LinkTargetFieldForNum(6)},
			validateReq: func(req *pb.GetLinkRequest) bool {
				return req.Name == "" &&
					req.SourceStage == "" &&
					req.SourceField == "" &&
					req.TargetStage == "" &&
					req.TargetField == testutil.LinkTargetFieldForNum(6)
			},
			responses: []*pb.Link{},
			output: [][]string{
				{
					NameText,
					SourceStageText,
					SourceFieldText,
					TargetStageText,
					TargetFieldText,
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
					LinkManagementServer: &ipb.MockLinkManagementServer{
						GetLinkFn: func(
							req *pb.GetLinkRequest,
							stream pb.LinkManagement_GetServer,
						) error {
							if !test.validateReq(req) {
								return fmt.Errorf(
									"validation failed with req %v",
									req,
								)
							}
							for _, l := range test.responses {
								if err := stream.Send(l); err != nil {
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
				cmd := NewCmdGetLink()
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

func pbLinkForNum(num int) *pb.Link {
	return &pb.Link{
		Name:        testutil.LinkNameForNumStr(num),
		SourceStage: testutil.LinkSourceStageForNumStr(num),
		SourceField: testutil.LinkSourceFieldForNum(num),
		TargetStage: testutil.LinkTargetStageForNumStr(num),
		TargetField: testutil.LinkTargetFieldForNum(num),
	}
}
