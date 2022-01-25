package reflection

import (
	"context"
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"github.com/DuarteMRAlves/maestro/internal/reflection"
	"sync"
)

type Manager struct {
	Rpcs sync.Map
}

func (m *Manager) FindRpc(
	ctx context.Context,
	name apitypes.StageName,
	query *reflection.FindQuery,
) error {
	panic("method not implemented")
}

func (m *Manager) GetRpc(stage apitypes.StageName) (reflection.RPC, bool) {
	rpc, ok := m.Rpcs.Load(stage)
	if !ok {
		return nil, false
	}
	return rpc.(reflection.RPC), true
}
