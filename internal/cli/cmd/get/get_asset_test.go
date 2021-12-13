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
	"regexp"
	"testing"
)

// TestGetAsset_CorrectDisplay performs integration testing on the GetAsset
// command considering operations that produce table outputs. It runs a maestro
// server and then executes a get asset command with predetermined arguments,
// verifying its output by comparing with an expected table.
func TestGetAsset_CorrectDisplay(t *testing.T) {
	tests := []struct {
		name        string
		defaultAddr bool
		args        []string
		assets      []*resources.AssetResource
		output      [][]string
	}{
		{
			name:   "empty assets",
			args:   []string{},
			assets: []*resources.AssetResource{},
			output: [][]string{
				{NameText, ImageText},
			},
		},
		{
			name:   "one asset",
			args:   []string{},
			assets: []*resources.AssetResource{assetForNum(0)},
			output: [][]string{
				{NameText, ImageText},
				{assetNameForNum(0), assetImageForNum(0)},
			},
		},
		{
			name: "multiple assets",
			args: []string{},
			assets: []*resources.AssetResource{
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
			name:        "multiple assets default address",
			defaultAddr: true,
			args:        []string{},
			assets: []*resources.AssetResource{
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
			assets: []*resources.AssetResource{
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
			assets: []*resources.AssetResource{
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
			assets: []*resources.AssetResource{
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
			assets: []*resources.AssetResource{
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

				if test.defaultAddr {
					lis = testutil.LockAndListenDefaultAddr(t)
					defer testutil.UnlockDefaultAddr()
				} else {
					lis = testutil.ListenAvailablePort(t)
				}

				addr = lis.Addr().String()

				if !test.defaultAddr {
					test.args = append(test.args, "--addr", addr)
				}

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

				populateAssets(t, test.assets, addr)

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

// TestGetAsset_CLIErrors performs integration testing on the GetAsset
// command with sets of flags that do no required the server to be running.
func TestGetAsset_CLIErrors(t *testing.T) {
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
				cmd := NewCmdGetAsset()
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
