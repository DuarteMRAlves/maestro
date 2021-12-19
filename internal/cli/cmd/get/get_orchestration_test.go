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

// TestGetOrchestration_CorrectDisplay performs integration testing on the
// GetOrchestration command considering operations that produce table outputs.
// It runs a maestro server and then executes a get orchestration command with
// predetermined arguments, verifying its output by comparing with an expected
// table.
func TestGetOrchestration_CorrectDisplay(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		orchestrations []*resources.OrchestrationSpec
		output         [][]string
	}{
		{
			name:           "empty orchestrations",
			args:           []string{},
			orchestrations: []*resources.OrchestrationSpec{},
			output:         [][]string{{NameText}},
		},
		{
			name: "one orchestration",
			args: []string{},
			orchestrations: []*resources.OrchestrationSpec{
				orchestrationForNum(0),
			},
			output: [][]string{
				{NameText}, {orchestrationNameForNum(0)},
			},
		},
		{
			name: "multiple orchestrations",
			args: []string{},
			orchestrations: []*resources.OrchestrationSpec{
				orchestrationForNum(1),
				orchestrationForNum(0),
				orchestrationForNum(2),
			},
			output: [][]string{
				{
					NameText,
				},
				{
					orchestrationNameForNum(0),
				},
				{
					orchestrationNameForNum(1),
				},
				{
					orchestrationNameForNum(2),
				},
			},
		},
		{
			name: "filter by name",
			args: []string{orchestrationNameForNum(2)},
			orchestrations: []*resources.OrchestrationSpec{
				orchestrationForNum(2),
				orchestrationForNum(1),
				orchestrationForNum(0),
			},
			output: [][]string{
				{
					NameText,
				},
				{
					orchestrationNameForNum(2),
				},
			},
		},
		{
			name: "no such name",
			args: []string{orchestrationNameForNum(3)},
			orchestrations: []*resources.OrchestrationSpec{
				orchestrationForNum(2),
				orchestrationForNum(1),
				orchestrationForNum(0),
			},
			output: [][]string{
				{
					NameText,
				},
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

				s, err := server.NewBuilder().
					WithGrpc().
					WithLogger(testutil.NewLogger(t)).
					Build()

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

				err = populateOrchestrations(t, test.orchestrations, addr)
				assert.NilError(t, err, "populate orchestrations")

				b := bytes.NewBufferString("")
				cmd := NewCmdGetOrchestration()
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
