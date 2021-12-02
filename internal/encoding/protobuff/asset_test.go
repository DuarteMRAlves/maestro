package protobuff

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/asset"
	"gotest.tools/v3/assert"
	"testing"
)

func TestMarshalAsset(t *testing.T) {
	tests := []struct {
		in asset.Asset
	}{
		{asset.Asset{Name: assetName, Image: assetImage}},
		{
			asset.Asset{
				Name:  assetName,
				Image: assetImage,
			},
		},
		{asset.Asset{Name: assetName}},
		{asset.Asset{Image: assetImage}},
	}

	for _, inner := range tests {
		in := inner.in
		name := fmt.Sprintf("in='%v'", in)

		t.Run(
			name, func(t *testing.T) {
				res := MarshalAsset(&in)
				assert.Equal(t, in.Name, res.Name, "Asset Name")
				assert.Equal(t, in.Image, res.Image, "Asset Image")
			})
	}
}

func TestUnmarshalAsset(t *testing.T) {
	tests := []struct {
		in *pb.Asset
	}{
		{&pb.Asset{Name: "Asset Name"}},
		{&pb.Asset{Name: "Asset Name"}},
	}

	for _, inner := range tests {
		in := inner.in
		name := fmt.Sprintf("in='%v'", in)

		t.Run(
			name, func(t *testing.T) {
				res, err := UnmarshalAsset(in)
				assert.Equal(t, nil, err, "Error")
				assert.Equal(t, in.Name, res.Name, "Asset Name")
			})
	}
}
