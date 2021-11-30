package create

import "github.com/DuarteMRAlves/maestro/internal/errdefs"

const (
	assetKind = "asset"
	stageKind = "stage"
)

func isValidKinds(resources []*Resource) error {
	for _, r := range resources {
		if !isValidKind(r.Kind) {
			return errdefs.InvalidArgumentWithMsg("invalid kind %v", r.Kind)
		}
	}
	return nil
}

func isValidKind(kind string) bool {
	return kind == assetKind || kind == stageKind
}

func isAssetKind(r *Resource) bool {
	return r.Kind == assetKind
}

func isStageKind(r *Resource) bool {
	return r.Kind == stageKind
}
