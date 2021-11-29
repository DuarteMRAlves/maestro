package asset

import (
	"github.com/DuarteMRAlves/maestro/internal/test"
	"testing"
)

func TestAsset_Clone(t *testing.T) {
	src := &Asset{Name: assetName, Image: assetImage}
	dst := src.Clone()
	test.DeepEqual(t, assetName, src.Name, "Source Name")
	test.DeepEqual(t, assetImage, src.Image, "Source Image")
	test.DeepEqual(t, assetName, dst.Name, "Dest Name")
	test.DeepEqual(t, assetImage, dst.Image, "Dest Image")
}
