package storage

import (
	"github.com/DuarteMRAlves/maestro/internal/api"
	"gotest.tools/v3/assert"
	"testing"
)

func TestStage_Clone(t *testing.T) {
	const (
		oldName    api.StageName = "Old Name"
		oldPhase                 = api.StageRunning
		oldAsset   api.AssetName = "Old Asset Name"
		oldAddress               = "OldAddress"
		oldService               = "OldService"
		oldRpc                   = "OldRpc"

		newName    api.StageName = "New Name"
		newPhase                 = api.StageFailed
		newAsset   api.AssetName = "New Asset Name"
		newAddress               = "NewAddress"
		newService               = "NewService"
		newRpc                   = "NewRpc"
	)

	s := &Stage{
		name:  oldName,
		phase: oldPhase,
		asset: oldAsset,
		rpcSpec: &RpcSpec{
			address: oldAddress,
			service: oldService,
			rpc:     oldRpc,
		},
	}
	c := s.Clone()
	assert.Equal(t, oldName, c.name, "cloned old Name")
	assert.Equal(t, oldPhase, c.phase, "cloned old phase")
	assert.Equal(t, oldAsset, c.asset, "cloned old asset id")
	assert.Equal(t, oldAddress, c.rpcSpec.address, "cloned old address")
	assert.Equal(t, oldService, c.rpcSpec.service, "cloned old service")
	assert.Equal(t, oldRpc, c.rpcSpec.rpc, "cloned old rpc")

	c.name = newName
	c.phase = newPhase
	c.asset = newAsset
	c.rpcSpec.address = newAddress
	c.rpcSpec.service = newService
	c.rpcSpec.rpc = newRpc

	assert.Equal(t, oldName, s.name, "source old Name")
	assert.Equal(t, oldPhase, s.phase, "source old phase")
	assert.Equal(t, oldAsset, s.asset, "source old asset id")
	assert.Equal(t, oldAddress, s.rpcSpec.address, "source old address")
	assert.Equal(t, oldService, s.rpcSpec.service, "source old service")
	assert.Equal(t, oldRpc, s.rpcSpec.rpc, "source old rpc")

	assert.Equal(t, newName, c.name, "cloned new Name")
	assert.Equal(t, newPhase, c.phase, "cloned new phase")
	assert.Equal(t, newAsset, c.asset, "cloned new asset id")
	assert.Equal(t, newAddress, c.rpcSpec.address, "cloned new address")
	assert.Equal(t, newService, c.rpcSpec.service, "cloned new service")
	assert.Equal(t, newRpc, c.rpcSpec.rpc, "cloned new rpc")
}
