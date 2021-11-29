package stage

import (
	"github.com/DuarteMRAlves/maestro/internal/test"
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
	test.DeepEqual(t, oldName, c.Name, "cloned old Name")
	test.DeepEqual(t, oldAsset, c.Asset, "cloned old asset id")
	test.DeepEqual(t, oldService, c.Service, "cloned old service")
	test.DeepEqual(t, oldMethod, c.Method, "cloned old method")

	c.Name = newName
	c.Asset = newAsset
	c.Service = newService
	c.Method = newMethod

	test.DeepEqual(t, oldName, s.Name, "source old Name")
	test.DeepEqual(t, oldAsset, s.Asset, "source old asset id")
	test.DeepEqual(t, oldService, s.Service, "source old service")
	test.DeepEqual(t, oldMethod, s.Method, "source old method")

	test.DeepEqual(t, newName, c.Name, "cloned new Name")
	test.DeepEqual(t, newAsset, c.Asset, "cloned new asset id")
	test.DeepEqual(t, newService, c.Service, "cloned new service")
	test.DeepEqual(t, newMethod, c.Method, "cloned new method")
}
