package get

import (
	"bytes"
	"github.com/DuarteMRAlves/maestro/internal/cli/resources"
	"github.com/DuarteMRAlves/maestro/internal/server"
	"github.com/pterm/pterm"
	"gotest.tools/v3/assert"
	"io/ioutil"
	"net"
	"regexp"
	"testing"
)

// TestGetLink_CorrectDisplay performs integration testing on the GetLink
// command considering operations that produce table outputs. It runs a maestro
// server and then executes a get link command with predetermined arguments,
// verifying its output by comparing with an expected table.
func TestGetLink_CorrectDisplay(t *testing.T) {
	tests := []struct {
		name   string
		args   []string
		links  []*resources.LinkResource
		output [][]string
	}{
		{
			name:  "empty links",
			args:  []string{},
			links: []*resources.LinkResource{},
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
			name:  "one link",
			args:  []string{},
			links: []*resources.LinkResource{linkForNum(0)},
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
			links: []*resources.LinkResource{
				linkForNum(1),
				linkForNum(0),
				linkForNum(2),
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
			links: []*resources.LinkResource{
				linkForNum(2),
				linkForNum(1),
				linkForNum(0),
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
			links: []*resources.LinkResource{
				linkForNum(1),
				linkForNum(2),
				linkForNum(0),
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
			links: []*resources.LinkResource{
				linkForNum(2),
				linkForNum(1),
				linkForNum(0),
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
			},
		},
		{
			name: "filter by target stage",
			args: []string{"--target-stage", linkTargetStageForNum(1)},
			links: []*resources.LinkResource{
				linkForNum(2),
				linkForNum(0),
				linkForNum(1),
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
			links: []*resources.LinkResource{
				linkForNum(0),
				linkForNum(2),
				linkForNum(1),
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
			links: []*resources.LinkResource{
				linkForNum(2),
				linkForNum(1),
				linkForNum(0),
			},
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
			links: []*resources.LinkResource{
				linkForNum(1),
				linkForNum(2),
				linkForNum(0),
			},
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
			links: []*resources.LinkResource{
				linkForNum(2),
				linkForNum(1),
				linkForNum(0),
			},
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
			links: []*resources.LinkResource{
				linkForNum(2),
				linkForNum(0),
				linkForNum(1),
			},
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
			links: []*resources.LinkResource{
				linkForNum(0),
				linkForNum(1),
				linkForNum(2),
			},
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
	assets := []*resources.AssetResource{
		assetForNum(0),
		assetForNum(1),
		assetForNum(2),
		assetForNum(3),
	}
	stages := []*resources.StageResource{
		stageForNum(0),
		stageForNum(1),
		stageForNum(2),
		stageForNum(3),
	}
	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				addr := "localhost:50051"
				lis, err := net.Listen("tcp", addr)
				assert.NilError(t, err, "failed to listen")
				s := server.NewBuilder().WithGrpc().Build()

				go func() {
					if err := s.ServeGrpc(lis); err != nil {
						t.Errorf("Failed to serve: %v", err)
						return
					}
				}()
				// Stop the server. Any calls in the test should be finished.
				// If not, an error should be raised.
				defer s.StopGrpc()

				populateAssets(t, assets, addr)
				populateStages(t, stages, addr)
				populateLinks(t, test.links, addr)

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
				assert.NilError(t, err, "render error")
				assert.Equal(t, expectedOut, string(out), "output differs")
			})
	}
}

// TestGetLink_CLIErrors performs integration testing on the GetLink
// command with sets of flags that do no required the server to be running.
func TestGetLink_CLIErrors(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectedOut string
	}{
		{
			"server not connected",
			[]string{},
			`unavailable: connection error: desc = "transport: Error while dialing dial tcp .+:50051: connect: connection refused"`,
		},
	}
	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				b := bytes.NewBufferString("")
				cmd := NewCmdGetLink()
				cmd.SetOut(b)
				cmd.SetArgs(test.args)
				err := cmd.Execute()
				assert.NilError(t, err, "execute error")
				out, err := ioutil.ReadAll(b)
				assert.NilError(t, err, "read output error")
				// This is not ideal but its to match the not connected error
				// with no ip. Detailed in GitHub issue
				// https://github.com/DuarteMRAlves/maestro/issues/29.
				matched, err := regexp.MatchString(
					test.expectedOut,
					string(out))
				assert.NilError(t, err, "matched output")
				assert.Assert(t, matched, "output not matched")
			})
	}
}
