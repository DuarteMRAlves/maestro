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

// TestGetLink_CorrectDisplay performs testing on the GetLink command
// considering operations that produce table outputs. It runs a mock maestro
// server and then executes a get link command with predetermined arguments,
// verifying its output by comparing with an expected table.
func TestGetLink_CorrectDisplay(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		validateQuery func(query *pb.Link) bool
		responses     []*pb.Link
		output        [][]string
	}{
		{
			name: "empty links",
			args: []string{},
			validateQuery: func(query *pb.Link) bool {
				return query.Name == "" &&
					query.SourceStage == "" &&
					query.SourceField == "" &&
					query.TargetStage == "" &&
					query.TargetField == ""
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
			validateQuery: func(query *pb.Link) bool {
				return query.Name == "" &&
					query.SourceStage == "" &&
					query.SourceField == "" &&
					query.TargetStage == "" &&
					query.TargetField == ""
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
					linkNameForNum(0),
					linkSourceStageForNum(0),
					linkSourceFieldForNum(0),
					linkTargetStageForNum(0),
					linkTargetFieldForNum(0),
				},
			},
		},
		{
			name: "multiple links",
			args: []string{},
			validateQuery: func(query *pb.Link) bool {
				return query.Name == "" &&
					query.SourceStage == "" &&
					query.SourceField == "" &&
					query.TargetStage == "" &&
					query.TargetField == ""
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
					linkNameForNum(0),
					linkSourceStageForNum(0),
					linkSourceFieldForNum(0),
					linkTargetStageForNum(0),
					linkTargetFieldForNum(0),
				},
				{
					linkNameForNum(1),
					linkSourceStageForNum(1),
					linkSourceFieldForNum(1),
					linkTargetStageForNum(1),
					linkTargetFieldForNum(1),
				},
				{
					linkNameForNum(2),
					linkSourceStageForNum(2),
					linkSourceFieldForNum(2),
					linkTargetStageForNum(2),
					linkTargetFieldForNum(2),
				},
			},
		},
		{
			name: "filter by name",
			args: []string{linkNameForNum(2)},
			validateQuery: func(query *pb.Link) bool {
				return query.Name == linkNameForNum(2) &&
					query.SourceStage == "" &&
					query.SourceField == "" &&
					query.TargetStage == "" &&
					query.TargetField == ""
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
					linkNameForNum(2),
					linkSourceStageForNum(2),
					linkSourceFieldForNum(2),
					linkTargetStageForNum(2),
					linkTargetFieldForNum(2),
				},
			},
		},
		{
			name: "filter by source stage",
			args: []string{"--source-stage", linkSourceStageForNum(2)},
			validateQuery: func(query *pb.Link) bool {
				return query.Name == "" &&
					query.SourceStage == linkSourceStageForNum(2) &&
					query.SourceField == "" &&
					query.TargetStage == "" &&
					query.TargetField == ""
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
					linkNameForNum(2),
					linkSourceStageForNum(2),
					linkSourceFieldForNum(2),
					linkTargetStageForNum(2),
					linkTargetFieldForNum(2),
				},
			},
		},
		{
			name: "filter by source field",
			args: []string{"--source-field", linkSourceFieldForNum(0)},
			validateQuery: func(query *pb.Link) bool {
				return query.Name == "" &&
					query.SourceStage == "" &&
					query.SourceField == linkSourceFieldForNum(0) &&
					query.TargetStage == "" &&
					query.TargetField == ""
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
					linkNameForNum(0),
					linkSourceStageForNum(0),
					linkSourceFieldForNum(0),
					linkTargetStageForNum(0),
					linkTargetFieldForNum(0),
				},
			},
		},
		{
			name: "filter by target stage",
			args: []string{"--target-stage", linkTargetStageForNum(1)},
			validateQuery: func(query *pb.Link) bool {
				return query.Name == "" &&
					query.SourceStage == "" &&
					query.SourceField == "" &&
					query.TargetStage == linkTargetStageForNum(1) &&
					query.TargetField == ""
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
					linkNameForNum(1),
					linkSourceStageForNum(1),
					linkSourceFieldForNum(1),
					linkTargetStageForNum(1),
					linkTargetFieldForNum(1),
				},
			},
		},
		{
			name: "filter by target field",
			args: []string{"--target-field", linkTargetFieldForNum(2)},
			validateQuery: func(query *pb.Link) bool {
				return query.Name == "" &&
					query.SourceStage == "" &&
					query.SourceField == "" &&
					query.TargetStage == "" &&
					query.TargetField == linkTargetFieldForNum(2)
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
					linkNameForNum(2),
					linkSourceStageForNum(2),
					linkSourceFieldForNum(2),
					linkTargetStageForNum(2),
					linkTargetFieldForNum(2),
				},
			},
		},
		{
			name: "no such name",
			args: []string{linkNameForNum(3)},
			validateQuery: func(query *pb.Link) bool {
				return query.Name == linkNameForNum(3) &&
					query.SourceStage == "" &&
					query.SourceField == "" &&
					query.TargetStage == "" &&
					query.TargetField == ""
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
			args: []string{"--source-stage", linkSourceStageForNum(3)},
			validateQuery: func(query *pb.Link) bool {
				return query.Name == "" &&
					query.SourceStage == linkSourceStageForNum(3) &&
					query.SourceField == "" &&
					query.TargetStage == "" &&
					query.TargetField == ""
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
			args: []string{"--source-field", linkSourceFieldForNum(4)},
			validateQuery: func(query *pb.Link) bool {
				return query.Name == "" &&
					query.SourceStage == "" &&
					query.SourceField == linkSourceFieldForNum(4) &&
					query.TargetStage == "" &&
					query.TargetField == ""
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
			args: []string{"--target-stage", linkTargetStageForNum(5)},
			validateQuery: func(query *pb.Link) bool {
				return query.Name == "" &&
					query.SourceStage == "" &&
					query.SourceField == "" &&
					query.TargetStage == linkTargetStageForNum(5) &&
					query.TargetField == ""
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
			args: []string{"--target-field", linkTargetFieldForNum(6)},
			validateQuery: func(query *pb.Link) bool {
				return query.Name == "" &&
					query.SourceStage == "" &&
					query.SourceField == "" &&
					query.TargetStage == "" &&
					query.TargetField == linkTargetFieldForNum(6)
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

				mockServer := mock.MaestroServer{
					LinkManagementServer: &mock.LinkManagementServer{
						GetLinkFn: func(
							query *pb.Link,
							stream pb.LinkManagement_GetServer,
						) error {
							if !test.validateQuery(query) {
								return fmt.Errorf(
									"validation failed with query %v",
									query)
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
			})
	}
}

func pbLinkForNum(num int) *pb.Link {
	return &pb.Link{
		Name:        linkNameForNum(num),
		SourceStage: linkSourceStageForNum(num),
		SourceField: linkSourceFieldForNum(num),
		TargetStage: linkTargetStageForNum(num),
		TargetField: linkTargetFieldForNum(num),
	}
}
