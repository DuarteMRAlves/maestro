package asset

import (
	"github.com/DuarteMRAlves/maestro/internal/assert"
	"github.com/DuarteMRAlves/maestro/internal/identifier"
	"testing"
)

var testId = identifier.Id{Val: "AssetTestId"}

const (
	name  = "asset name"
	image = "asset image"
)

func TestAsset_Clone(t *testing.T) {
	src := &Asset{Id: testId, Name: name, Image: image}
	dst := src.Clone()
	assert.DeepEqual(t, testId, src.Id, "Source Id")
	assert.DeepEqual(t, name, src.Name, "Source Name")
	assert.DeepEqual(t, image, src.Image, "Source Image")
	assert.DeepEqual(t, testId, dst.Id, "Dest Id")
	assert.DeepEqual(t, name, dst.Name, "Dest Name")
	assert.DeepEqual(t, image, dst.Image, "Dest Image")
}
