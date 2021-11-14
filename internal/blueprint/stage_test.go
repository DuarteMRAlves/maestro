package blueprint

import (
	"github.com/DuarteMRAlves/maestro/internal/assert"
	"github.com/DuarteMRAlves/maestro/internal/identifier"
	"testing"
)

const (
	oldId      = "OALD23SJ25"
	oldName    = "Old Name"
	oldAssetId = "WF89FEF9WH"
	oldService = "OldService"
	oldMethod  = "OldMethod"

	newId      = "DG341FEF04"
	newName    = "New Name"
	newAssetId = "TREGE9878F"
	newService = "NewService"
	newMethod  = "NewMethod"
)

func TestStage_Clone(t *testing.T) {
	s := &Stage{
		Id:      identifier.Id{Val: oldId},
		Name:    oldName,
		AssetId: identifier.Id{Val: oldAssetId},
		Service: oldService,
		Method:  oldMethod,
	}
	c := s.Clone()
	assert.DeepEqual(t, oldId, c.Id.Val, "cloned old id")
	assert.DeepEqual(t, oldName, c.Name, "cloned old name")
	assert.DeepEqual(t, oldAssetId, c.AssetId.Val, "cloned old asset id")
	assert.DeepEqual(t, oldService, c.Service, "cloned old service")
	assert.DeepEqual(t, oldMethod, c.Method, "cloned old method")

	c.Id.Val = newId
	c.Name = newName
	c.AssetId.Val = newAssetId
	c.Service = newService
	c.Method = newMethod

	assert.DeepEqual(t, oldId, s.Id.Val, "source old id")
	assert.DeepEqual(t, oldName, s.Name, "source old name")
	assert.DeepEqual(t, oldAssetId, s.AssetId.Val, "source old asset id")
	assert.DeepEqual(t, oldService, s.Service, "source old service")
	assert.DeepEqual(t, oldMethod, s.Method, "source old method")

	assert.DeepEqual(t, newId, c.Id.Val, "cloned new id")
	assert.DeepEqual(t, newName, c.Name, "cloned new name")
	assert.DeepEqual(t, newAssetId, c.AssetId.Val, "cloned new asset id")
	assert.DeepEqual(t, newService, c.Service, "cloned new service")
	assert.DeepEqual(t, newMethod, c.Method, "cloned new method")
}
