package stage

import (
	testing2 "github.com/DuarteMRAlves/maestro/internal/testing"
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
	testing2.DeepEqual(t, oldName, c.Name, "cloned old Name")
	testing2.DeepEqual(t, oldAsset, c.Asset, "cloned old asset id")
	testing2.DeepEqual(t, oldService, c.Service, "cloned old service")
	testing2.DeepEqual(t, oldMethod, c.Method, "cloned old method")

	c.Name = newName
	c.Asset = newAsset
	c.Service = newService
	c.Method = newMethod

	testing2.DeepEqual(t, oldName, s.Name, "source old Name")
	testing2.DeepEqual(t, oldAsset, s.Asset, "source old asset id")
	testing2.DeepEqual(t, oldService, s.Service, "source old service")
	testing2.DeepEqual(t, oldMethod, s.Method, "source old method")

	testing2.DeepEqual(t, newName, c.Name, "cloned new Name")
	testing2.DeepEqual(t, newAsset, c.Asset, "cloned new asset id")
	testing2.DeepEqual(t, newService, c.Service, "cloned new service")
	testing2.DeepEqual(t, newMethod, c.Method, "cloned new method")
}
