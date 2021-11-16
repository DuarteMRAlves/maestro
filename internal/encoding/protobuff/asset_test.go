package protobuff

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/assert"
	"github.com/DuarteMRAlves/maestro/internal/asset"
	"github.com/DuarteMRAlves/maestro/internal/identifier"
	"testing"
)

func TestMarshalAsset(t *testing.T) {
	rand, _ := identifier.Rand(5)
	tests := []struct {
		in asset.Asset
	}{
		{asset.Asset{Id: rand, Name: assetName, Image: assetImage}},
		{
			asset.Asset{
				Id:    identifier.Empty(),
				Name:  assetName,
				Image: assetImage,
			},
		},
		{asset.Asset{Name: assetName}},
		{asset.Asset{Id: rand}},
		{asset.Asset{Image: assetImage}},
	}

	for _, inner := range tests {
		in := inner.in
		name := fmt.Sprintf("in='%v'", in)

		t.Run(
			name, func(t *testing.T) {
				res := MarshalAsset(&in)
				assert.DeepEqual(t, in.Id.Val, res.Id.Val, "Asset Id")
				assert.DeepEqual(t, in.Name, res.Name, "Asset Name")
				assert.DeepEqual(t, in.Image, res.Image, "Asset Image")
			})
	}
}

func TestUnmarshalAsset(t *testing.T) {
	tests := []struct {
		in *pb.Asset
	}{
		{&pb.Asset{Id: &pb.Id{Val: "Some String"}, Name: "Asset Name"}},
		{&pb.Asset{Id: &pb.Id{Val: ""}, Name: "Asset Name"}},
		{&pb.Asset{Name: "Asset Name"}},
		{&pb.Asset{Id: &pb.Id{Val: "Some String"}}},
	}

	for _, inner := range tests {
		in := inner.in
		name := fmt.Sprintf("in='%v'", in)

		t.Run(
			name, func(t *testing.T) {
				res, err := UnmarshalAsset(in)
				assert.DeepEqual(t, nil, err, "Error")
				if in.Id != nil {
					assert.DeepEqual(t, in.Id.Val, res.Id.Val, "Asset Id")
				} else {
					assert.DeepEqual(
						t,
						identifier.Empty(),
						res.Id,
						"Asset Empty Id")
				}
				assert.DeepEqual(t, in.Name, res.Name, "Asset Name")
			})
	}
}
