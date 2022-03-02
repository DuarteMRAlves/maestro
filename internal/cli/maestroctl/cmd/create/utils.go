package create

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

func collectAssetSpecs(resources []*Resource) ([]*AssetSpec, error) {
	requests := make([]*AssetSpec, 0)
	for _, r := range resources {
		if r.IsAssetKind() {
			req, ok := r.Spec.(*AssetSpec)
			if !ok {
				return nil, errdefs.InternalWithMsg(
					"create asset request spec cast failed: %v",
					r,
				)
			}
			requests = append(requests, req)
		}
	}
	return requests, nil
}

func collectStageSpecs(resources []*Resource) ([]*StageSpec, error) {
	requests := make([]*StageSpec, 0)
	for _, r := range resources {
		if r.IsStageKind() {
			req, ok := r.Spec.(*StageSpec)
			if !ok {
				return nil, errdefs.InternalWithMsg(
					"stage spec cast failed: %v",
					r,
				)
			}
			requests = append(requests, req)
		}
	}
	return requests, nil
}

func collectLinkSpecs(resources []*Resource) ([]*LinkSpec, error) {
	requests := make([]*LinkSpec, 0)
	for _, r := range resources {
		if r.IsLinkKind() {
			req, ok := r.Spec.(*LinkSpec)
			if !ok {
				return nil, errdefs.InternalWithMsg(
					"link spec cast failed: %v",
					r,
				)
			}
			requests = append(requests, req)
		}
	}
	return requests, nil
}

func collectOrchestrationSpecs(resources []*Resource) (
	[]*OrchestrationSpec,
	error,
) {
	requests := make([]*OrchestrationSpec, 0)
	for _, r := range resources {
		if r.IsOrchestrationKind() {
			req, ok := r.Spec.(*OrchestrationSpec)
			if !ok {
				return nil, errdefs.InternalWithMsg(
					"orchestration spec cast failed> %v",
					req,
				)
			}
			requests = append(requests, req)
		}
	}
	return requests, nil
}
