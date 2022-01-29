package storage

import (
	"github.com/DuarteMRAlves/maestro/internal/api"
)

func buildOrchestrationQueryFilter(
	req *api.GetOrchestrationRequest,
) func(b *api.Orchestration) bool {
	filters := make([]func(b *api.Orchestration) bool, 0)
	if req.Name != "" {
		filters = append(
			filters,
			func(b *api.Orchestration) bool {
				return b.Name == req.Name
			},
		)
	}
	if req.Phase != "" {
		filters = append(
			filters,
			func(o *api.Orchestration) bool {
				return o.Phase == req.Phase
			},
		)
	}
	switch len(filters) {
	case 0:
		return func(b *api.Orchestration) bool { return true }
	case 1:
		return filters[0]
	default:
		return func(b *api.Orchestration) bool {
			for _, f := range filters {
				if !f(b) {
					return false
				}
			}
			return true
		}
	}
}

func buildStageQueryFilter(req *api.GetStageRequest) func(s *api.Stage) bool {
	filters := make([]func(s *api.Stage) bool, 0)
	if req.Name != "" {
		filters = append(
			filters,
			func(s *api.Stage) bool {
				return s.Name == req.Name
			},
		)
	}
	if req.Phase != "" {
		filters = append(
			filters,
			func(s *api.Stage) bool {
				return s.Phase == req.Phase
			},
		)
	}
	if req.Asset != "" {
		filters = append(
			filters,
			func(s *api.Stage) bool {
				return s.Asset == req.Asset
			},
		)
	}
	if req.Service != "" {
		filters = append(
			filters,
			func(s *api.Stage) bool {
				return s.Service == req.Service
			},
		)
	}
	if req.Rpc != "" {
		filters = append(
			filters,
			func(s *api.Stage) bool {
				return s.Rpc == req.Rpc
			},
		)
	}
	if req.Address != "" {
		filters = append(
			filters,
			func(s *api.Stage) bool {
				return s.Address == req.Address
			},
		)
	}
	if req.Orchestration != "" {
		filters = append(
			filters,
			func(s *api.Stage) bool {
				return s.Orchestration == req.Orchestration
			},
		)
	}
	switch len(filters) {
	case 0:
		return func(s *api.Stage) bool { return true }
	case 1:
		return filters[0]
	default:
		return func(s *api.Stage) bool {
			for _, f := range filters {
				if !f(s) {
					return false
				}
			}
			return true
		}
	}
}

func buildLinkQueryFilter(req *api.GetLinkRequest) func(l *api.Link) bool {
	filters := make([]func(l *api.Link) bool, 0)
	if req.Name != "" {
		filters = append(
			filters,
			func(l *api.Link) bool {
				return l.Name == req.Name
			},
		)
	}
	if req.SourceStage != "" {
		filters = append(
			filters,
			func(l *api.Link) bool {
				return l.SourceStage == req.SourceStage
			},
		)
	}
	if req.SourceField != "" {
		filters = append(
			filters,
			func(l *api.Link) bool {
				return l.SourceField == req.SourceField
			},
		)
	}
	if req.TargetStage != "" {
		filters = append(
			filters,
			func(l *api.Link) bool {
				return l.TargetStage == req.TargetStage
			},
		)
	}
	if req.TargetField != "" {
		filters = append(
			filters,
			func(l *api.Link) bool {
				return l.TargetField == req.TargetField
			},
		)
	}
	if req.Orchestration != "" {
		filters = append(
			filters,
			func(l *api.Link) bool {
				return l.Orchestration == req.Orchestration
			},
		)
	}
	switch len(filters) {
	case 0:
		return func(l *api.Link) bool { return true }
	case 1:
		return filters[0]
	default:
		return func(l *api.Link) bool {
			for _, f := range filters {
				if !f(l) {
					return false
				}
			}
			return true
		}
	}
}

func buildAssetQueryFilter(req *api.GetAssetRequest) func(a *api.Asset) bool {
	filters := make([]func(a *api.Asset) bool, 0)
	if req.Name != "" {
		filters = append(
			filters,
			func(a *api.Asset) bool {
				return a.Name == req.Name
			},
		)
	}
	if req.Image != "" {
		filters = append(
			filters,
			func(a *api.Asset) bool {
				return a.Image == req.Image
			},
		)
	}
	if len(filters) > 0 {
		return func(a *api.Asset) bool {
			for _, f := range filters {
				if !f(a) {
					return false
				}
			}
			return true
		}
	}
	return func(a *api.Asset) bool {
		return true
	}
}
