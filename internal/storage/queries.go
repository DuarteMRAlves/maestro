package storage

import apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"

func buildOrchestrationQueryFilter(
	query *apitypes.Orchestration,
) func(b *Orchestration) bool {
	filters := make([]func(b *Orchestration) bool, 0)
	if query.Name != "" {
		filters = append(
			filters,
			func(b *Orchestration) bool {
				return b.Name() == query.Name
			},
		)
	}
	if query.Phase != "" {
		filters = append(
			filters,
			func(o *Orchestration) bool {
				return o.phase == query.Phase
			},
		)
	}
	switch len(filters) {
	case 0:
		return func(b *Orchestration) bool { return true }
	case 1:
		return filters[0]
	default:
		return func(b *Orchestration) bool {
			for _, f := range filters {
				if !f(b) {
					return false
				}
			}
			return true
		}
	}
}

func buildStageQueryFilter(query *apitypes.Stage) func(s *Stage) bool {
	filters := make([]func(s *Stage) bool, 0)
	if query.Name != "" {
		filters = append(
			filters,
			func(s *Stage) bool {
				return s.name == query.Name
			},
		)
	}
	if query.Phase != "" {
		filters = append(
			filters,
			func(s *Stage) bool {
				return s.phase == query.Phase
			},
		)
	}
	if query.Asset != "" {
		filters = append(
			filters,
			func(s *Stage) bool {
				return s.asset == query.Asset
			},
		)
	}
	if query.Service != "" {
		filters = append(
			filters,
			func(s *Stage) bool {
				return s.rpcSpec.service == query.Service
			},
		)
	}
	if query.Rpc != "" {
		filters = append(
			filters,
			func(s *Stage) bool {
				return s.rpcSpec.rpc == query.Rpc
			},
		)
	}
	if query.Address != "" {
		filters = append(
			filters,
			func(s *Stage) bool {
				return s.rpcSpec.address == query.Address
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

func buildLinkQueryFilter(query *apitypes.Link) func(l *Link) bool {
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

func buildQueryFilter(query *apitypes.Asset) func(a *Asset) bool {
	filters := make([]func(a *Asset) bool, 0)
	if query.Name != "" {
		filters = append(
			filters,
			func(a *Asset) bool {
				return a.Name() == query.Name
			},
		)
	}
	if query.Image != "" {
		filters = append(
			filters,
			func(a *Asset) bool {
				return a.Image() == query.Image
			},
		)
	}
	if len(filters) > 0 {
		return func(a *Asset) bool {
			for _, f := range filters {
				if !f(a) {
					return false
				}
			}
			return true
		}
	}
	return func(a *Asset) bool {
		return true
	}
}
