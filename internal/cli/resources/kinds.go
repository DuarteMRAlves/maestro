package resources

import (
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
)

const (
	AssetKind = "asset"
	StageKind = "stage"
	LinkKind  = "link"
)

func IsValidKinds(resources []*Resource) error {
	for _, r := range resources {
		if !IsValidKind(r.Kind) {
			return errdefs.InvalidArgumentWithMsg("invalid kind '%v'", r.Kind)
		}
	}
	return nil
}

func IsValidKind(kind string) bool {
	return kind == AssetKind || kind == StageKind || kind == LinkKind
}

func IsAssetKind(r *Resource) bool {
	return r.Kind == AssetKind
}

func IsStageKind(r *Resource) bool {
	return r.Kind == StageKind
}

func IsLinkKind(r *Resource) bool {
	return r.Kind == LinkKind
}

// FilterAssets creates a new array with the resources whose kind is asset.
func FilterAssets(resources []*Resource) []*Resource {
	assets := make([]*Resource, 0)
	for _, r := range resources {
		if IsAssetKind(r) {
			assets = append(assets, r)
		}
	}
	return assets
}

// FilterStages creates a new array with the resources whose kind is stage.
func FilterStages(resources []*Resource) []*Resource {
	stages := make([]*Resource, 0)
	for _, r := range resources {
		if IsStageKind(r) {
			stages = append(stages, r)
		}
	}
	return stages
}

// FilterLinks creates a new array with the resources whose kind is asset.
func FilterLinks(resources []*Resource) []*Resource {
	links := make([]*Resource, 0)
	for _, r := range resources {
		if IsLinkKind(r) {
			links = append(links, r)
		}
	}
	return links
}
