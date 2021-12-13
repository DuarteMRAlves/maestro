package create

import (
	"bytes"
	"github.com/DuarteMRAlves/maestro/internal/server"
	"github.com/DuarteMRAlves/maestro/internal/testutil"
	"gotest.tools/v3/assert"
	"io/ioutil"
	"net"
	"regexp"
	"testing"
)

// TestCreateWithServer performs integration testing on the Create command
// considering operations that require the server to be running.
// It runs a maestro server and executes the command with predetermined
// arguments, verifying its output.
func TestCreateWithServer(t *testing.T) {
	tests := []struct {
		name        string
		defaultAddr bool
		args        []string
		expectedOut string
	}{
		{
			"multiple resources in a single file",
			false,
			[]string{"-f", "../../../../tests/resources/create/resources.yml"},
			"",
		},
		{
			"multiple resources in multiple files",
			false,
			[]string{
				"-f",
				"../../../../tests/resources/create/stages.yml",
				"-f",
				"../../../../tests/resources/create/links.yml",
				"-f",
				"../../../../tests/resources/create/assets.yml",
			},
			"",
		},
		{
			"custom address",
			true,
			[]string{"-f", "../../../../tests/resources/create/resources.yml"},
			"",
		},
		{
			"asset not found",
			false,
			[]string{
				"-f",
				"../../../../tests/resources/create/asset_not_found.yml",
			},
			"not found: asset 'unknown-asset' not found",
		},
		{
			"stage not found",
			false,
			[]string{
				"-f",
				"../../../../tests/resources/create/stage_not_found.yml",
			},
			"not found: target stage 'unknown-stage' not found",
		},
	}
	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
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

				s := server.NewBuilder().WithGrpc().Build()

				go func() {
					if err := s.ServeGrpc(lis); err != nil {
						t.Errorf("Failed to serve: %v", err)
						return
					}
				}()
				defer func() {
					// Stop the server. Any calls in the test should be finished.
					// If not, an error should be raised.
					s.StopGrpc()
				}()

				b := bytes.NewBufferString("")
				cmd := NewCmdCreate()
				cmd.SetOut(b)
				cmd.SetArgs(test.args)
				err = cmd.Execute()
				assert.NilError(t, err, "execute error")
				out, err := ioutil.ReadAll(b)
				assert.NilError(t, err, "read output error")
				assert.Equal(t, test.expectedOut, string(out), "output differs")
			})
	}
}

// TestCreateWithServer performs integration testing on the Create command
// considering operations that do not require the server to be running.
func TestCreateWithoutServer(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectedOut string
	}{
		{
			"no files",
			[]string{},
			"invalid argument: please specify input files",
		},
		{
			"no such file",
			[]string{"-f", "missing_file.yml"},
			"invalid argument: open missing_file.yml: no such file or directory",
		},
		{
			"invalid kind",
			[]string{
				"-f",
				"../../../../tests/resources/create/invalid_kind.yml",
			},
			"invalid argument: invalid kind 'invalid-kind'",
		},
		{
			"invalid specs",
			[]string{
				"-f",
				"../../../../tests/resources/create/invalid_specs.yml",
			},
			"invalid argument: unknown spec fields: invalid_spec_1,invalid_spec_2",
		},
		{
			"server not connected",
			[]string{"-f", "../../../../tests/resources/create/resources.yml"},
			`unavailable: connection error: desc = "transport: Error while dialing dial tcp .+:50051: connect: connection refused"`,
		},
	}
	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				b := bytes.NewBufferString("")
				cmd := NewCmdCreate()
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
