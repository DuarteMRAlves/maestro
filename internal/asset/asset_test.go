package asset

import (
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"gotest.tools/v3/assert"
	"testing"
)

func TestAsset_Clone(t *testing.T) {
	const (
		assetName  apitypes.AssetName = "Asset-Name"
		assetImage                    = "Asset-Image"
	)
	src := &Asset{Name: assetName, Image: assetImage}
	dst := src.Clone()
	assert.Equal(t, assetName, src.Name, "source name")
	assert.Equal(t, assetImage, src.Image, "source image")
	assert.Equal(t, assetName, dst.Name, "dest name")
	assert.Equal(t, assetImage, dst.Image, "dest image")
}
