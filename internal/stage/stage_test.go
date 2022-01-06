package stage

import (
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"gotest.tools/v3/assert"
	"testing"
)

func TestStage_Clone(t *testing.T) {
	const (
		oldName    = "Old Name"
		oldPhase   = apitypes.StageRunning
		oldAsset   = "Old Asset Name"
		oldAddress = "OldAddress"

		newName    = "New Name"
		newPhase   = apitypes.StageFailed
		newAsset   = "New Asset Name"
		newAddress = "NewAddress"
	)
	s := &Stage{
		Name:    oldName,
		Phase:   oldPhase,
		Asset:   oldAsset,
		Address: oldAddress,
	}
	c := s.Clone()
	assert.Equal(t, oldName, c.Name, "cloned old Name")
	assert.Equal(t, oldPhase, c.Phase, "cloned old phase")
	assert.Equal(t, oldAsset, c.Asset, "cloned old asset id")
	assert.Equal(t, oldAddress, c.Address, "cloned old address")

	c.Name = newName
	c.Phase = newPhase
	c.Asset = newAsset
	c.Address = newAddress

	assert.Equal(t, oldName, s.Name, "source old Name")
	assert.Equal(t, oldPhase, s.Phase, "source old phase")
	assert.Equal(t, oldAsset, s.Asset, "source old asset id")
	assert.Equal(t, oldAddress, s.Address, "source old address")

	assert.Equal(t, newName, c.Name, "cloned new Name")
	assert.Equal(t, newPhase, c.Phase, "cloned new phase")
	assert.Equal(t, newAsset, c.Asset, "cloned new asset id")
	assert.Equal(t, newAddress, c.Address, "cloned new address")
}
