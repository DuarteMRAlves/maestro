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

func buildStageQueryFilter(req *api.GetStageRequest) func(s *Stage) bool {
	filters := make([]func(s *Stage) bool, 0)
	if req.Name != "" {
		filters = append(
			filters,
			func(s *Stage) bool {
				return s.name == req.Name
			},
		)
	}
	if req.Phase != "" {
		filters = append(
			filters,
			func(s *Stage) bool {
				return s.phase == req.Phase
			},
		)
	}
	if req.Asset != "" {
		filters = append(
			filters,
			func(s *Stage) bool {
				return s.asset == req.Asset
			},
		)
	}
	if req.Service != "" {
		filters = append(
			filters,
			func(s *Stage) bool {
				return s.rpcSpec.service == req.Service
			},
		)
	}
	if req.Rpc != "" {
		filters = append(
			filters,
			func(s *Stage) bool {
				return s.rpcSpec.rpc == req.Rpc
			},
		)
	}
	if req.Address != "" {
		filters = append(
			filters,
			func(s *Stage) bool {
				return s.rpcSpec.address == req.Address
			},
		)
	}
	switch len(filters) {
	case 0:
		return func(s *Stage) bool { return true }
	case 1:
		return filters[0]
	default:
		return func(s *Stage) bool {
			for _, f := range filters {
				if !f(s) {
					return false
				}
			}
			return true
		}
	}
}

func buildLinkQueryFilter(query *api.Link) func(l *Link) bool {
	filters := make([]func(l *Link) bool, 0)
	if query.Name != "" {
		filters = append(
			filters,
			func(l *Link) bool {
				return l.name == query.Name
			},
		)
	}
	if query.SourceStage != "" {
		filters = append(
			filters,
			func(l *Link) bool {
				return l.sourceStage == query.SourceStage
			},
		)
	}
	if query.SourceField != "" {
		filters = append(
			filters,
			func(l *Link) bool {
				return l.sourceField == query.SourceField
			},
		)
	}
	if query.TargetStage != "" {
		filters = append(
			filters,
			func(l *Link) bool {
				return l.targetStage == query.TargetStage
			},
		)
	}
	if query.TargetField != "" {
		filters = append(
			filters,
			func(l *Link) bool {
				return l.targetField == query.TargetField
			},
		)
	}
	switch len(filters) {
	case 0:
		return func(l *Link) bool { return true }
	case 1:
		return filters[0]
	default:
		return func(l *Link) bool {
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