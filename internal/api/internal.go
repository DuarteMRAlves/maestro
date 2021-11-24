package api

import (
	"github.com/DuarteMRAlves/maestro/internal/asset"
	"github.com/DuarteMRAlves/maestro/internal/blueprint"
	"github.com/DuarteMRAlves/maestro/internal/link"
	"github.com/DuarteMRAlves/maestro/internal/stage"
)

// InternalAPI is an interface that collects all the available commands
// for the maestro server. All calls on external APIs should be redirected
// through this API that collects all functionality.
type InternalAPI interface {
	CreateAsset(config *asset.Asset) error
	GetAsset(query *asset.Asset) []*asset.Asset

	CreateStage(config *stage.Stage) error
	GetStage(query *stage.Stage) []*stage.Stage

	CreateLink(config *link.Link) error
	GetLink(query *link.Link) []*link.Link

	CreateBlueprint(config *blueprint.Blueprint) error
}
