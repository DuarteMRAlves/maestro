package create

import "github.com/DuarteMRAlves/maestro/internal/errdefs"

const (
	assetKind = "asset"
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
	return kind == assetKind
}

func isAssetKind(r *Resource) bool {
	return r.Kind == assetKind
}
