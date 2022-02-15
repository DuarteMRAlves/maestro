package get

import (
	"bytes"
	"fmt"
	"github.com/DuarteMRAlves/maestro/api/pb"
	ipb "github.com/DuarteMRAlves/maestro/internal/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/util"
	"github.com/pterm/pterm"
	"gotest.tools/v3/assert"
	"io/ioutil"
	"net"
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
					Name:  util.AssetNameForNumStr(0),
					Image: util.AssetImageForNum(0),
				},
			},
			output: [][]string{
				{NameText, ImageText},
				{util.AssetNameForNumStr(0), util.AssetImageForNum(0)},
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
					Name:  util.AssetNameForNumStr(2),
					Image: util.AssetImageForNum(2),
				},
				{
					Name:  util.AssetNameForNumStr(1),
					Image: util.AssetImageForNum(1),
				},
				{
					Name:  util.AssetNameForNumStr(0),
					Image: util.AssetImageForNum(0),
				},
			},
			output: [][]string{
				{NameText, ImageText},
				{util.AssetNameForNumStr(0), util.AssetImageForNum(0)},
				{util.AssetNameForNumStr(1), util.AssetImageForNum(1)},
				{util.AssetNameForNumStr(2), util.AssetImageForNum(2)},
			},
		},
		{
			name: "filter by name",
			args: []string{util.AssetNameForNumStr(1)},
			validateQuery: func(req *pb.GetAssetRequest) bool {
				return req.Name == util.AssetNameForNumStr(1) &&
					req.Image == ""
			},
			responses: []*pb.Asset{
				{
					Name:  util.AssetNameForNumStr(1),
					Image: util.AssetImageForNum(1),
				},
			},
			output: [][]string{
				{NameText, ImageText},
				{util.AssetNameForNumStr(1), util.AssetImageForNum(1)},
			},
		},
		{
			name: "filter by image",
			args: []string{"--image", util.AssetImageForNum(2)},
			validateQuery: func(req *pb.GetAssetRequest) bool {
				return req.Name == "" &&
					req.Image == util.AssetImageForNum(2)
			},
			responses: []*pb.Asset{
				{
					Name:  util.AssetNameForNumStr(2),
					Image: util.AssetImageForNum(2),
				},
			},
			output: [][]string{
				{NameText, ImageText},
				{util.AssetNameForNumStr(2), util.AssetImageForNum(2)},
			},
		},
		{
			name: "no such name",
			args: []string{util.AssetNameForNumStr(3)},
			validateQuery: func(req *pb.GetAssetRequest) bool {
				return req.Name == util.AssetNameForNumStr(3) &&
					req.Image == ""
			},
			responses: []*pb.Asset{},
			output: [][]string{
				{NameText, ImageText},
			},
		},
		{
			name: "no such image",
			args: []string{"--image", util.AssetImageForNum(4)},
			validateQuery: func(req *pb.GetAssetRequest) bool {
				return req.Name == "" &&
					req.Image == util.AssetImageForNum(4)
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
				lis, err := net.Listen("tcp", "localhost:0")
				assert.NilError(t, err, "failed to listen")

				addr := lis.Addr().String()
				test.args = append(test.args, "--maestro", addr)

				mockServer := ipb.MockMaestroServer{
					ArchitectureManagementServer: &ipb.MockArchitectureManagementServer{
						GetAssetFn: func(
							req *pb.GetAssetRequest,
							stream pb.ArchitectureManagement_GetAssetServer,
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
			},
		)
	}
}
