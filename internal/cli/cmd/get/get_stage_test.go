package get

import (
	"bytes"
	"github.com/DuarteMRAlves/maestro/internal/cli/resources"
	"github.com/DuarteMRAlves/maestro/internal/server"
	"github.com/DuarteMRAlves/maestro/internal/testutil"
	"github.com/pterm/pterm"
	"gotest.tools/v3/assert"
	"io/ioutil"
	"net"
	"testing"
)

// TestGetStage_CorrectDisplay performs integration testing on the GetStage
// command considering operations that produce table outputs. It runs a maestro
// server and then executes a get stage command with predetermined arguments,
// verifying its output by comparing with an expected table.
func TestGetStage_CorrectDisplay(t *testing.T) {
	tests := []struct {
		name   string
		args   []string
		stages []*resources.StageSpec
		output [][]string
	}{
		{
			name:   "empty stages",
			args:   []string{},
			stages: []*resources.StageSpec{},
			output: [][]string{
				{NameText, AssetText, ServiceText, MethodText, AddressText},
			},
		},
		{
			name:   "one stage",
			args:   []string{},
			stages: []*resources.StageSpec{stageForNum(0)},
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
			stages: []*resources.StageSpec{
				stageForNum(1),
				stageForNum(0),
				stageForNum(2),
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
			stages: []*resources.StageSpec{
				stageForNum(2),
				stageForNum(1),
				stageForNum(0),
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
			stages: []*resources.StageSpec{
				stageForNum(1),
				stageForNum(2),
				stageForNum(0),
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
			stages: []*resources.StageSpec{
				stageForNum(2),
				stageForNum(1),
				stageForNum(0),
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
			stages: []*resources.StageSpec{
				stageForNum(2),
				stageForNum(0),
				stageForNum(1),
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
			stages: []*resources.StageSpec{
				stageForNum(2),
				stageForNum(1),
				stageForNum(0),
			},
			output: [][]string{
				{NameText, AssetText, ServiceText, MethodText, AddressText},
			},
		},
		{
			name: "no such asset",
			args: []string{"--asset", assetNameForNum(3)},
			stages: []*resources.StageSpec{
				stageForNum(1),
				stageForNum(2),
				stageForNum(0),
			},
			output: [][]string{
				{NameText, AssetText, ServiceText, MethodText, AddressText},
			},
		},
		{
			name: "no such service",
			args: []string{"--service", stageServiceForNum(4)},
			stages: []*resources.StageSpec{
				stageForNum(2),
				stageForNum(1),
				stageForNum(0),
			},
			output: [][]string{
				{NameText, AssetText, ServiceText, MethodText, AddressText},
			},
		},
		{
			name: "no such method",
			args: []string{"--method", stageMethodForNum(5)},
			stages: []*resources.StageSpec{
				stageForNum(2),
				stageForNum(0),
				stageForNum(1),
			},
			output: [][]string{
				{NameText, AssetText, ServiceText, MethodText, AddressText},
			},
		},
	}
	assets := []*resources.AssetSpec{
		assetForNum(0),
		assetForNum(1),
		assetForNum(2),
	}
	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				var (
					lis  net.Listener
					addr string
					err  error
				)

				lis = testutil.ListenAvailablePort(t)

				addr = lis.Addr().String()

				test.args = append(test.args, "--addr", addr)

				s, err := server.NewBuilder().WithGrpc().Build()
				assert.NilError(t, err, "build server")

				go func() {
					if err := s.ServeGrpc(lis); err != nil {
						t.Errorf("Failed to serve: %v", err)
						return
					}
				}()
				// Stop the server. Any calls in the test should be finished.
				// If not, an error should be raised.
				defer s.StopGrpc()

				err = populateAssets(t, assets, addr)
				assert.NilError(t, err, "populate assets")
				err = populateStages(t, test.stages, addr)
				assert.NilError(t, err, "populate stages")

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
			})
	}
}
