package create

import (
	"bytes"
	"github.com/DuarteMRAlves/maestro/internal/cli/cmd/create"
	"github.com/DuarteMRAlves/maestro/internal/server"
	"gotest.tools/v3/assert"
	"io/ioutil"
	"net"
	"testing"
)

// TestCreate performs integration testing on the Create command
// It runs a maestro server and executes the command with predetermined
// arguments, verifying its output.
func TestCreate(t *testing.T) {
	tests := []struct {
		name        string
		serverAddr  string
		args        []string
		expectedOut string
	}{
		{
			"multiple resources in a single file",
			"localhost:50051",
			[]string{"-f", "../resources/resources.yml"},
			"",
		},
		{
			"multiple resources in multiple files",
			"localhost:50051",
			[]string{
				"-f",
				"../resources/stages.yml",
				"-f",
				"../resources/links.yml",
				"-f",
				"../resources/assets.yml",
			},
			"",
		},
		{
			"custom address",
			"localhost:50052",
			[]string{
				"-f",
				"../resources/resources.yml",
				"--addr",
				"localhost:50052",
			},
			"",
		},
		{
			"no files",
			"localhost:50051",
			[]string{},
			"invalid argument: please specify input files",
		},
	}
	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				lis, err := net.Listen("tcp", test.serverAddr)
				assert.NilError(t, err, "failed to listen")

				s := server.NewBuilder().WithGrpc().Build()

				go func() {
					if err := s.ServeGrpc(lis); err != nil {
						t.Fatalf("Failed to serve: %v", err)
					}
				}()
				defer s.GracefulStopGrpc()

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
