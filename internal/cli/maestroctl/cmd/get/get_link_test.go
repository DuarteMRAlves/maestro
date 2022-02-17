package get

import (
	"bytes"
	"fmt"
	"github.com/DuarteMRAlves/maestro/api/pb"
	ipb "github.com/DuarteMRAlves/maestro/internal/api/pb"
	"github.com/pterm/pterm"
	"gotest.tools/v3/assert"
	"io/ioutil"
	"net"
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
					"link-0",
					"stage-0",
					"source-field-0",
					"stage-1",
					"target-field-0",
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
					"link-0",
					"stage-0",
					"source-field-0",
					"stage-1",
					"target-field-0",
				},
				{
					"link-1",
					"stage-1",
					"source-field-1",
					"stage-2",
					"target-field-1",
				},
				{
					"link-2",
					"stage-2",
					"source-field-2",
					"stage-3",
					"target-field-2",
				},
			},
		},
		{
			name: "filter by name",
			args: []string{"link-2"},
			validateReq: func(req *pb.GetLinkRequest) bool {
				return req.Name == "link-2" &&
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
					"link-2",
					"stage-2",
					"source-field-2",
					"stage-3",
					"target-field-2",
				},
			},
		},
		{
			name: "filter by source stage",
			args: []string{
				"--source-stage",
				"stage-2",
			},
			validateReq: func(req *pb.GetLinkRequest) bool {
				return req.Name == "" &&
					req.SourceStage == "stage-2" &&
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
					"link-2",
					"stage-2",
					"source-field-2",
					"stage-3",
					"target-field-2",
				},
			},
		},
		{
			name: "filter by source field",
			args: []string{"--source-field", "source-field-0"},
			validateReq: func(req *pb.GetLinkRequest) bool {
				return req.Name == "" &&
					req.SourceStage == "" &&
					req.SourceField == "source-field-0" &&
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
					"link-0",
					"stage-0",
					"source-field-0",
					"stage-1",
					"target-field-0",
				},
			},
		},
		{
			name: "filter by target stage",
			args: []string{
				"--target-stage",
				"stage-2",
			},
			validateReq: func(req *pb.GetLinkRequest) bool {
				return req.Name == "" &&
					req.SourceStage == "" &&
					req.SourceField == "" &&
					req.TargetStage == "stage-2" &&
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
					"link-1",
					"stage-1",
					"source-field-1",
					"stage-2",
					"target-field-1",
				},
			},
		},
		{
			name: "filter by target field",
			args: []string{"--target-field", "target-field-2"},
			validateReq: func(req *pb.GetLinkRequest) bool {
				return req.Name == "" &&
					req.SourceStage == "" &&
					req.SourceField == "" &&
					req.TargetStage == "" &&
					req.TargetField == "target-field-2"
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
					"link-2",
					"stage-2",
					"source-field-2",
					"stage-3",
					"target-field-2",
				},
			},
		},
		{
			name: "no such name",
			args: []string{"link-3"},
			validateReq: func(req *pb.GetLinkRequest) bool {
				return req.Name == "link-3" &&
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
			args: []string{"--source-stage", "stage-3"},
			validateReq: func(req *pb.GetLinkRequest) bool {
				return req.Name == "" &&
					req.SourceStage == "stage-3" &&
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
			args: []string{"--source-field", "source-field-4"},
			validateReq: func(req *pb.GetLinkRequest) bool {
				return req.Name == "" &&
					req.SourceStage == "" &&
					req.SourceField == "source-field-4" &&
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
			args: []string{"--target-stage", "stage-6"},
			validateReq: func(req *pb.GetLinkRequest) bool {
				return req.Name == "" &&
					req.SourceStage == "" &&
					req.SourceField == "" &&
					req.TargetStage == "stage-6" &&
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
			args: []string{"--target-field", "target-field-6"},
			validateReq: func(req *pb.GetLinkRequest) bool {
				return req.Name == "" &&
					req.SourceStage == "" &&
					req.SourceField == "" &&
					req.TargetStage == "" &&
					req.TargetField == "target-field-6"
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
				lis, err := net.Listen("tcp", "localhost:0")
				assert.NilError(t, err, "failed to listen")

				addr := lis.Addr().String()
				test.args = append(test.args, "--maestro", addr)

				mockServer := ipb.MockMaestroServer{
					ArchitectureManagementServer: &ipb.MockArchitectureManagementServer{
						GetLinkFn: func(
							req *pb.GetLinkRequest,
							stream pb.ArchitectureManagement_GetLinkServer,
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

func pbLinkForNum(num int) *pb.Link {
	return &pb.Link{
		Name:        fmt.Sprintf("link-%d", num),
		SourceStage: fmt.Sprintf("stage-%d", num),
		SourceField: fmt.Sprintf("source-field-%d", num),
		TargetStage: fmt.Sprintf("stage-%d", num+1),
		TargetField: fmt.Sprintf("target-field-%d", num),
	}
}
