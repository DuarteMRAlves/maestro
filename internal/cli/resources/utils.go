package resources

import (
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
)

func IsValidKinds(resources []*Resource) error {
	for _, r := range resources {
		if !r.IsValidKind() {
			return errdefs.InvalidArgumentWithMsg("invalid kind '%v'", r.Kind)
		}
	}
	return nil
}

// FilterAssets creates a new array with the resources whose kind is asset.
func FilterAssets(resources []*Resource) []*Resource {
	assets := make([]*Resource, 0)
	for _, r := range resources {
		if r.IsAssetKind() {
			assets = append(assets, r)
		}
	}
	return assets
}

// FilterStages creates a new array with the resources whose kind is stage.
func FilterStages(resources []*Resource) []*Resource {
	stages := make([]*Resource, 0)
	for _, r := range resources {
		if r.IsStageKind() {
			stages = append(stages, r)
		}
	}
	return stages
}

// FilterLinks creates a new array with the resources whose kind is asset.
func FilterLinks(resources []*Resource) []*Resource {
	links := make([]*Resource, 0)
	for _, r := range resources {
		if r.IsLinkKind() {
			links = append(links, r)
		}
	}
	return links
}

// FilterOrchestrations creates a new array with the resources whose kind is
// orchestration.
func FilterOrchestrations(resources []*Resource) []*Resource {
	orchestrations := make([]*Resource, 0)
	for _, r := range resources {
		if r.IsOrchestrationKind() {
			orchestrations = append(orchestrations, r)
		}
	}
	return orchestrations
}
