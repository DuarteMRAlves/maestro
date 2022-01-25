package api

import (
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
)

// InternalAPI is an interface that collects all the available commands
// for the maestro server. All calls on external APIs should be redirected
// through this API that collects all functionality.
type InternalAPI interface {
	CreateAsset(config *apitypes.Asset) error
	GetAsset(query *apitypes.Asset) ([]*apitypes.Asset, error)

	CreateStage(config *apitypes.Stage) error
	GetStage(query *apitypes.Stage) ([]*apitypes.Stage, error)

	CreateLink(config *apitypes.Link) error
	GetLink(query *apitypes.Link) []*apitypes.Link

	CreateOrchestration(config *apitypes.Orchestration) error
	GetOrchestration(
		query *apitypes.Orchestration,
	) ([]*apitypes.Orchestration, error)
}
