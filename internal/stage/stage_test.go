package stage

import (
	"gotest.tools/v3/assert"
	"testing"
)

const (
	oldName    = "Old Name"
	oldAsset   = "Old Asset Name"
	oldAddress = "OldAddress"

	newName    = "New Name"
	newAsset   = "New Asset Name"
	newAddress = "NewAddress"
)

func TestStage_Clone(t *testing.T) {
	s := &Stage{
		Name:    oldName,
		Asset:   oldAsset,
		Address: oldAddress,
	}
	c := s.Clone()
	assert.Equal(t, oldName, c.Name, "cloned old Name")
	assert.Equal(t, oldAsset, c.Asset, "cloned old asset id")
	assert.Equal(t, oldAddress, c.Address, "cloned old address")

	c.Name = newName
	c.Asset = newAsset
	c.Address = newAddress

	assert.Equal(t, oldName, s.Name, "source old Name")
	assert.Equal(t, oldAsset, s.Asset, "source old asset id")
	assert.Equal(t, oldAddress, s.Address, "source old address")

	assert.Equal(t, newName, c.Name, "cloned new Name")
	assert.Equal(t, newAsset, c.Asset, "cloned new asset id")
	assert.Equal(t, newAddress, c.Address, "cloned new address")
}
