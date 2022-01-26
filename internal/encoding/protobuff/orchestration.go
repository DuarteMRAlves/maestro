package protobuff

import (
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/validate"
)

// MarshalOrchestration returns a protobuf encoding of the given orchestration.
func MarshalOrchestration(o *api.Orchestration) (
	*pb.Orchestration,
	error,
) {
	if ok, err := validate.ArgNotNil(o, "o"); !ok {
		return nil, err
	}
	links := make([]string, 0, len(o.Links))
	for _, l := range o.Links {
		links = append(links, string(l))
	}
	protoBp := &pb.Orchestration{
		Name:  string(o.Name),
		Phase: string(o.Phase),
		Links: links,
	}
	return protoBp, nil
}

// UnmarshalOrchestration returns an orchestration from the orchestration
// protobuf encoding.
func UnmarshalOrchestration(p *pb.Orchestration) (
	*api.Orchestration,
	error,
) {
	if ok, err := validate.ArgNotNil(p, "p"); !ok {
		return nil, err
	}

	links := make([]api.LinkName, 0, len(p.Links))
	for _, l := range p.Links {
		links = append(links, api.LinkName(l))
	}

	return &api.Orchestration{
		Name:  api.OrchestrationName(p.Name),
		Phase: api.OrchestrationPhase(p.Phase),
		Links: links,
	}, nil
}
