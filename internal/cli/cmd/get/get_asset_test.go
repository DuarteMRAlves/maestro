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

// TestGetAsset_CorrectDisplay performs integration testing on the GetAsset
// command considering operations that produce table outputs. It runs a maestro
// server and then executes a get asset command with predetermined arguments,
// verifying its output by comparing with an expected table.
func TestGetAsset_CorrectDisplay(t *testing.T) {
	tests := []struct {
		name   string
		args   []string
		assets []*resources.AssetSpec
		output [][]string
	}{
		{
			name:   "empty assets",
			args:   []string{},
			assets: []*resources.AssetSpec{},
			output: [][]string{
				{NameText, ImageText},
			},
		},
		{
			name:   "one asset",
			args:   []string{},
			assets: []*resources.AssetSpec{assetForNum(0)},
			output: [][]string{
				{NameText, ImageText},
				{assetNameForNum(0), assetImageForNum(0)},
			},
		},
		{
			name: "multiple assets",
			args: []string{},
			assets: []*resources.AssetSpec{
				assetForNum(0),
				assetForNum(2),
				assetForNum(1),
			},
			output: [][]string{
				{NameText, ImageText},
				{assetNameForNum(0), assetImageForNum(0)},
				{assetNameForNum(1), assetImageForNum(1)},
				{assetNameForNum(2), assetImageForNum(2)},
			},
		},
		{
			name: "filter by name",
			args: []string{assetNameForNum(1)},
			assets: []*resources.AssetSpec{
				assetForNum(2),
				assetForNum(0),
				assetForNum(1),
			},
			output: [][]string{
				{NameText, ImageText},
				{assetNameForNum(1), assetImageForNum(1)},
			},
		},
		{
			name: "filter by image",
			args: []string{"--image", assetImageForNum(2)},
			assets: []*resources.AssetSpec{
				assetForNum(1),
				assetForNum(0),
				assetForNum(2),
			},
			output: [][]string{
				{NameText, ImageText},
				{assetNameForNum(2), assetImageForNum(2)},
			},
		},
		{
			name: "no such name",
			args: []string{assetNameForNum(3)},
			assets: []*resources.AssetSpec{
				assetForNum(2),
				assetForNum(0),
				assetForNum(1),
			},
			output: [][]string{
				{NameText, ImageText},
			},
		},
		{
			name: "no such image",
			args: []string{"--image", assetImageForNum(4)},
			assets: []*resources.AssetSpec{
				assetForNum(1),
				assetForNum(0),
				assetForNum(2),
			},
			output: [][]string{
				{NameText, ImageText},
			},
		},
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

				err = populateAssets(t, test.assets, addr)
				assert.NilError(t, err, "populate assets")

				b := bytes.NewBufferString("")
				cmd := NewCmdGetAsset()
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
