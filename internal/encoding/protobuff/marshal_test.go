package protobuff

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/assert"
	"github.com/DuarteMRAlves/maestro/internal/asset"
	"github.com/DuarteMRAlves/maestro/internal/identifier"
	"testing"
)

const (
	name  = "Asset Name"
	image = "user/image:version"
)

func TestMarshalID(t *testing.T) {
	rand, _ := identifier.Rand(5)
	tests := []struct {
		in identifier.Id
	}{
		{identifier.Empty()},
		{rand},
	}
	for _, inner := range tests {
		in := inner.in
		name := fmt.Sprintf("in='%v'", in)

		t.Run(
			name, func(t *testing.T) {
				res := MarshalID(in)
				assert.DeepEqual(t, res.Val, in.Val, "Id Val")
			})
	}
}

func TestMarshalAsset(t *testing.T) {
	rand, _ := identifier.Rand(5)
	tests := []struct {
		in asset.Asset
	}{
		{asset.Asset{Id: rand, Name: name, Image: image}},
		{asset.Asset{Id: identifier.Empty(), Name: name, Image: image}},
		{asset.Asset{Name: name}},
		{asset.Asset{Id: rand}},
		{asset.Asset{Image: image}},
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
