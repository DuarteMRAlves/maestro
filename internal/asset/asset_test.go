package asset

import (
	testing2 "github.com/DuarteMRAlves/maestro/internal/testing"
	"testing"
)

func TestAsset_Clone(t *testing.T) {
	src := &Asset{Name: assetName, Image: assetImage}
	dst := src.Clone()
	testing2.DeepEqual(t, assetName, src.Name, "Source Name")
	testing2.DeepEqual(t, assetImage, src.Image, "Source Image")
	testing2.DeepEqual(t, assetName, dst.Name, "Dest Name")
	testing2.DeepEqual(t, assetImage, dst.Image, "Dest Image")
}
