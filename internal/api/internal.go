package api

import (
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
)

// InternalAPI is an interface that collects all the available commands
// for the maestro server. All calls on external APIs should be redirected
// through this API that collects all functionality.
type InternalAPI interface {
	CreateAsset(config *CreateAssetRequest) error
	GetAsset(query *GetAssetRequest) ([]*apitypes.Asset, error)

	CreateStage(config *apitypes.Stage) error
	GetStage(query *apitypes.Stage) ([]*apitypes.Stage, error)

	CreateLink(config *apitypes.Link) error
	GetLink(query *apitypes.Link) ([]*apitypes.Link, error)

	CreateOrchestration(config *apitypes.Orchestration) error
	GetOrchestration(
		query *apitypes.Orchestration,
	) ([]*apitypes.Orchestration, error)
}
