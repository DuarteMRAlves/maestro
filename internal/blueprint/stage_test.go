package blueprint

import (
	"github.com/DuarteMRAlves/maestro/internal/assert"
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
	assert.DeepEqual(t, oldName, c.Name, "cloned old name")
	assert.DeepEqual(t, oldAsset, c.Asset, "cloned old asset id")
	assert.DeepEqual(t, oldService, c.Service, "cloned old service")
	assert.DeepEqual(t, oldMethod, c.Method, "cloned old method")

	c.Name = newName
	c.Asset = newAsset
	c.Service = newService
	c.Method = newMethod

	assert.DeepEqual(t, oldName, s.Name, "source old name")
	assert.DeepEqual(t, oldAsset, s.Asset, "source old asset id")
	assert.DeepEqual(t, oldService, s.Service, "source old service")
	assert.DeepEqual(t, oldMethod, s.Method, "source old method")

	assert.DeepEqual(t, newName, c.Name, "cloned new name")
	assert.DeepEqual(t, newAsset, c.Asset, "cloned new asset id")
	assert.DeepEqual(t, newService, c.Service, "cloned new service")
	assert.DeepEqual(t, newMethod, c.Method, "cloned new method")
}
