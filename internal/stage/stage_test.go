package stage

import (
	"gotest.tools/v3/assert"
	"testing"
)

const (
	oldName    = "Old Name"
	oldAsset   = "Old Asset Name"
	oldService = "OldService"
	oldMethod  = "OldMethod"

	newName    = "New Name"
	newAsset   = "New Asset Name"
	newService = "NewService"
	newMethod  = "NewMethod"
)

func TestStage_Clone(t *testing.T) {
	s := &Stage{
		Name:    oldName,
		Asset:   oldAsset,
		Service: oldService,
		Method:  oldMethod,
	}
	c := s.Clone()
	assert.Equal(t, oldName, c.Name, "cloned old Name")
	assert.Equal(t, oldAsset, c.Asset, "cloned old asset id")
	assert.Equal(t, oldService, c.Service, "cloned old service")
	assert.Equal(t, oldMethod, c.Method, "cloned old method")

	c.Name = newName
	c.Asset = newAsset
	c.Service = newService
	c.Method = newMethod

	assert.Equal(t, oldName, s.Name, "source old Name")
	assert.Equal(t, oldAsset, s.Asset, "source old asset id")
	assert.Equal(t, oldService, s.Service, "source old service")
	assert.Equal(t, oldMethod, s.Method, "source old method")

	assert.Equal(t, newName, c.Name, "cloned new Name")
	assert.Equal(t, newAsset, c.Asset, "cloned new asset id")
	assert.Equal(t, newService, c.Service, "cloned new service")
	assert.Equal(t, newMethod, c.Method, "cloned new method")
}
