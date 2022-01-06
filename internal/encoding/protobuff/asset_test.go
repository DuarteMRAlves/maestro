package protobuff

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/api/pb"
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"gotest.tools/v3/assert"
	"testing"
)

func TestMarshalAsset(t *testing.T) {
	const (
		name  apitypes.AssetName = "Asset Name"
		image                    = "user/image:version"
	)
	tests := []struct {
		in apitypes.Asset
	}{
		{apitypes.Asset{Name: name, Image: image}},
		{
			apitypes.Asset{
				Name:  name,
				Image: image,
			},
		},
		{apitypes.Asset{Name: name}},
		{apitypes.Asset{Image: image}},
	}

	for _, inner := range tests {
		in := inner.in
		name := fmt.Sprintf("in='%v'", in)

		t.Run(
			name, func(t *testing.T) {
				res, err := MarshalAsset(&in)
				assert.NilError(t, err, "marshal error")
				assert.Equal(t, string(in.Name), res.Name, "Asset Name")
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
				assert.Equal(t, in.Name, string(res.Name), "Asset Name")
			})
	}
}

func TestMarshalAssetNil(t *testing.T) {
	res, err := MarshalAsset(nil)
	assert.Assert(t, errdefs.IsInvalidArgument(err), "err type")
	assert.ErrorContains(t, err, "'a' is nil")
	assert.Assert(t, res == nil, "nil return value")
}

func TestUnmarshalAssetNil(t *testing.T) {
	res, err := UnmarshalAsset(nil)
	assert.Assert(t, errdefs.IsInvalidArgument(err), "err type")
	assert.ErrorContains(t, err, "'p' is nil")
	assert.Assert(t, res == nil, "nil return value")
}
