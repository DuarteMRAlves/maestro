package protobuff

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/asset"
	"github.com/DuarteMRAlves/maestro/internal/test"
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
				test.DeepEqual(t, in.Name, res.Name, "Asset Name")
				test.DeepEqual(t, in.Image, res.Image, "Asset Image")
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
				test.DeepEqual(t, nil, err, "Error")
				test.DeepEqual(t, in.Name, res.Name, "Asset Name")
			})
	}
}
