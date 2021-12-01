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
			return errdefs.InvalidArgumentWithMsg("invalid kind %v", r.Kind)
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
