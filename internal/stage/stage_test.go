package stage

import (
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"gotest.tools/v3/assert"
	"testing"
)

func TestStage_Clone(t *testing.T) {
	const (
		oldName    apitypes.StageName = "Old Name"
		oldPhase                      = apitypes.StageRunning
		oldAsset   apitypes.AssetName = "Old Asset Name"
		oldAddress                    = "OldAddress"

		newName    apitypes.StageName = "New Name"
		newPhase                      = apitypes.StageFailed
		newAsset   apitypes.AssetName = "New Asset Name"
		newAddress                    = "NewAddress"
	)
	s := &Stage{
		name:    oldName,
		phase:   oldPhase,
		asset:   oldAsset,
		address: oldAddress,
	}
	c := s.Clone()
	assert.Equal(t, oldName, c.name, "cloned old Name")
	assert.Equal(t, oldPhase, c.phase, "cloned old phase")
	assert.Equal(t, oldAsset, c.asset, "cloned old asset id")
	assert.Equal(t, oldAddress, c.address, "cloned old address")

	c.name = newName
	c.phase = newPhase
	c.asset = newAsset
	c.address = newAddress

	assert.Equal(t, oldName, s.name, "source old Name")
	assert.Equal(t, oldPhase, s.phase, "source old phase")
	assert.Equal(t, oldAsset, s.asset, "source old asset id")
	assert.Equal(t, oldAddress, s.address, "source old address")

	assert.Equal(t, newName, c.name, "cloned new Name")
	assert.Equal(t, newPhase, c.phase, "cloned new phase")
	assert.Equal(t, newAsset, c.asset, "cloned new asset id")
	assert.Equal(t, newAddress, c.address, "cloned new address")
}
