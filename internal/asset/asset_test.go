package asset

import (
	"github.com/DuarteMRAlves/maestro/internal/assert"
	"testing"
)

func TestAsset_Clone(t *testing.T) {
	src := &Asset{Name: assetName, Image: assetImage}
	dst := src.Clone()
	assert.DeepEqual(t, assetName, src.Name, "Source Name")
	assert.DeepEqual(t, assetImage, src.Image, "Source Image")
	assert.DeepEqual(t, assetName, dst.Name, "Dest Name")
	assert.DeepEqual(t, assetImage, dst.Image, "Dest Image")
}
