package get

import (
	"bytes"
	"fmt"
	"github.com/DuarteMRAlves/maestro/api/pb"
	ipb "github.com/DuarteMRAlves/maestro/internal/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/testutil"
	"github.com/pterm/pterm"
	"gotest.tools/v3/assert"
	"io/ioutil"
	"testing"
)

// TestGetAsset_CorrectDisplay performs testing on the GetAsset command
// considering operations that produce table outputs. It runs a mock maestro
// server and then executes a get asset command with predetermined arguments,
// verifying its output by comparing with an expected table.
func TestGetAsset_CorrectDisplay(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		validateQuery func(*pb.GetAssetRequest) bool
		responses     []*pb.Asset
		output        [][]string
	}{
		{
			name: "empty assets",
			args: []string{},
			validateQuery: func(req *pb.GetAssetRequest) bool {
				return req.Name == "" && req.Image == ""
			},
			responses: []*pb.Asset{},
			output: [][]string{
				{NameText, ImageText},
			},
		},
		{
			name: "one asset",
			args: []string{},
			validateQuery: func(req *pb.GetAssetRequest) bool {
				return req.Name == "" && req.Image == ""
			},
			responses: []*pb.Asset{
				{
					Name:  testutil.AssetNameForNumStr(0),
					Image: testutil.AssetImageForNum(0),
				},
			},
			output: [][]string{
				{NameText, ImageText},
				{testutil.AssetNameForNumStr(0), testutil.AssetImageForNum(0)},
			},
		},
		{
			name: "multiple assets",
			args: []string{},
			validateQuery: func(req *pb.GetAssetRequest) bool {
				return req.Name == "" && req.Image == ""
			},
			responses: []*pb.Asset{
				{
					Name:  testutil.AssetNameForNumStr(2),
					Image: testutil.AssetImageForNum(2),
				},
				{
					Name:  testutil.AssetNameForNumStr(1),
					Image: testutil.AssetImageForNum(1),
				},
				{
					Name:  testutil.AssetNameForNumStr(0),
					Image: testutil.AssetImageForNum(0),
				},
			},
			output: [][]string{
				{NameText, ImageText},
				{testutil.AssetNameForNumStr(0), testutil.AssetImageForNum(0)},
				{testutil.AssetNameForNumStr(1), testutil.AssetImageForNum(1)},
				{testutil.AssetNameForNumStr(2), testutil.AssetImageForNum(2)},
			},
		},
		{
			name: "filter by name",
			args: []string{testutil.AssetNameForNumStr(1)},
			validateQuery: func(req *pb.GetAssetRequest) bool {
				return req.Name == testutil.AssetNameForNumStr(1) &&
					req.Image == ""
			},
			responses: []*pb.Asset{
				{
					Name:  testutil.AssetNameForNumStr(1),
					Image: testutil.AssetImageForNum(1),
				},
			},
			output: [][]string{
				{NameText, ImageText},
				{testutil.AssetNameForNumStr(1), testutil.AssetImageForNum(1)},
			},
		},
		{
			name: "filter by image",
			args: []string{"--image", testutil.AssetImageForNum(2)},
			validateQuery: func(req *pb.GetAssetRequest) bool {
				return req.Name == "" &&
					req.Image == testutil.AssetImageForNum(2)
			},
			responses: []*pb.Asset{
				{
					Name:  testutil.AssetNameForNumStr(2),
					Image: testutil.AssetImageForNum(2),
				},
			},
			output: [][]string{
				{NameText, ImageText},
				{testutil.AssetNameForNumStr(2), testutil.AssetImageForNum(2)},
			},
		},
		{
			name: "no such name",
			args: []string{testutil.AssetNameForNumStr(3)},
			validateQuery: func(req *pb.GetAssetRequest) bool {
				return req.Name == testutil.AssetNameForNumStr(3) &&
					req.Image == ""
			},
			responses: []*pb.Asset{},
			output: [][]string{
				{NameText, ImageText},
			},
		},
		{
			name: "no such image",
			args: []string{"--image", testutil.AssetImageForNum(4)},
			validateQuery: func(req *pb.GetAssetRequest) bool {
				return req.Name == "" &&
					req.Image == testutil.AssetImageForNum(4)
			},
			responses: []*pb.Asset{},
			output: [][]string{
				{NameText, ImageText},
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

				mockServer := ipb.MockMaestroServer{
					AssetManagementServer: &ipb.MockAssetManagementServer{
						GetAssetFn: func(
							req *pb.GetAssetRequest,
							stream pb.AssetManagement_GetServer,
						) error {
							if !test.validateQuery(req) {
								return fmt.Errorf(
									"validation failed with req %v",
									req,
								)
							}
							for _, a := range test.responses {
								if err := stream.Send(a); err != nil {
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
				cmd := NewCmdGetAsset()
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
			},
		)
	}
}
