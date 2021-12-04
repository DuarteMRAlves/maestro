package create

import (
	"bytes"
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/cli/cmd/create"
	"github.com/DuarteMRAlves/maestro/internal/server"
	"gotest.tools/v3/assert"
	"io/ioutil"
	"net"
	"testing"
	"time"
)

// TestCreateWithServer performs integration testing on the Create command
// considering operations that require the server to be running.
// It runs a maestro server and executes the command with predetermined
// arguments, verifying its output.
func TestCreateWithServer(t *testing.T) {
	tests := []struct {
		name        string
		serverAddr  string
		args        []string
		expectedOut string
	}{
		{
			"multiple resources in a single file",
			"localhost:50051",
			[]string{"-f", "../resources/create/resources.yml"},
			"",
		},
		{
			"multiple resources in multiple files",
			"localhost:50051",
			[]string{
				"-f",
				"../resources/create/stages.yml",
				"-f",
				"../resources/create/links.yml",
				"-f",
				"../resources/create/assets.yml",
			},
			"",
		},
		{
			"custom address",
			"localhost:50052",
			[]string{
				"-f",
				"../resources/create/resources.yml",
				"--addr",
				"localhost:50052",
			},
			"",
		},
		{
			"asset not found",
			"localhost:50051",
			[]string{"-f", "../resources/create/asset_not_found.yml"},
			"not found: asset 'unknown-asset' not found",
		},
		{
			"stage not found",
			"localhost:50051",
			[]string{"-f", "../resources/create/stage_not_found.yml"},
			"not found: target stage 'unknown-stage' not found",
		},
	}
	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				fmt.Println(test.name)
				lis, err := net.Listen("tcp", test.serverAddr)
				assert.NilError(t, err, "failed to listen")
				s := server.NewBuilder().WithGrpc().Build()

				go func() {
					if err := s.ServeGrpc(lis); err != nil {
						fmt.Println(err)
						t.Errorf("Failed to serve: %v", err)
						return
					}
				}()
				defer func() {
					s.GracefulStopGrpc()
					// Wait a bit before forcefully stopping every call
					time.Sleep(10 * time.Millisecond)
					s.StopGrpc()
				}()

				b := bytes.NewBufferString("")
				cmd := create.NewCmdCreate()
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
			[]string{"-f", "../resources/create/invalid_kind.yml"},
			"invalid argument: invalid kind 'invalid-kind'",
		},
		{
			"invalid specs",
			[]string{"-f", "../resources/create/invalid_specs.yml"},
			"invalid argument: unknown spec fields: invalid_spec_1,invalid_spec_2",
		},
	}
	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				fmt.Println(test.name)

				b := bytes.NewBufferString("")
				cmd := create.NewCmdCreate()
				cmd.SetOut(b)
				cmd.SetArgs(test.args)
				err := cmd.Execute()
				assert.NilError(t, err, "execute error")
				out, err := ioutil.ReadAll(b)
				assert.NilError(t, err, "read output error")
				assert.Equal(t, test.expectedOut, string(out), "output differs")
			})
	}
}
