package resources

import (
	"github.com/DuarteMRAlves/maestro/internal/api"
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

// FilterCreateAssetRequests creates a new array with the requests to create
// assets.
func FilterCreateAssetRequests(resources []*Resource) (
	[]*api.CreateAssetRequest,
	error,
) {
	requests := make([]*api.CreateAssetRequest, 0)
	for _, r := range resources {
		if r.IsAssetKind() {
			req, ok := r.Spec.(*api.CreateAssetRequest)
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

// FilterCreateStageRequests creates a new array with the requests to create
// stages.
func FilterCreateStageRequests(resources []*Resource) (
	[]*api.CreateStageRequest,
	error,
) {
	requests := make([]*api.CreateStageRequest, 0)
	for _, r := range resources {
		if r.IsStageKind() {
			req, ok := r.Spec.(*api.CreateStageRequest)
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

// FilterCreateLinkRequests creates a new array with the requests to create
// links.
func FilterCreateLinkRequests(resources []*Resource) (
	[]*api.CreateLinkRequest,
	error,
) {
	requests := make([]*api.CreateLinkRequest, 0)
	for _, r := range resources {
		if r.IsLinkKind() {
			req, ok := r.Spec.(*api.CreateLinkRequest)
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

// FilterCreateOrchestrationRequests creates a new array with the requests to
// create orchestrations.
func FilterCreateOrchestrationRequests(resources []*Resource) (
	[]*api.CreateOrchestrationRequest,
	error,
) {
	requests := make([]*api.CreateOrchestrationRequest, 0)
	for _, r := range resources {
		if r.IsOrchestrationKind() {
			req, ok := r.Spec.(*api.CreateOrchestrationRequest)
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
