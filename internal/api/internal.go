package api

import (
	"github.com/DuarteMRAlves/maestro/internal/api/types"
	"github.com/DuarteMRAlves/maestro/internal/asset"
	"github.com/DuarteMRAlves/maestro/internal/link"
	"github.com/DuarteMRAlves/maestro/internal/orchestration"
)

// InternalAPI is an interface that collects all the available commands
// for the maestro server. All calls on external APIs should be redirected
// through this API that collects all functionality.
type InternalAPI interface {
	CreateAsset(config *asset.Asset) error
	GetAsset(query *asset.Asset) []*asset.Asset

	CreateStage(config *types.Stage) error
	GetStage(query *types.Stage) []*types.Stage

	CreateLink(config *link.Link) error
	GetLink(query *link.Link) []*link.Link

	CreateOrchestration(config *orchestration.Orchestration) error
	GetOrchestration(query *orchestration.Orchestration) []*orchestration.Orchestration
}
