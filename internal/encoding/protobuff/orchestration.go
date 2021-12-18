package protobuff

import (
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/orchestration"
	"github.com/DuarteMRAlves/maestro/internal/validate"
)

// MarshalOrchestration returns a protobuf encoding of the given orchestration.
func MarshalOrchestration(o *orchestration.Orchestration) (
	*pb.Orchestration,
	error,
) {
	if ok, err := validate.ArgNotNil(o, "o"); !ok {
		return nil, err
	}
	links := make([]string, 0, len(o.Links))
	for _, l := range o.Links {
		links = append(links, l)
	}
	protoBp := &pb.Orchestration{
		Name:  o.Name,
		Links: links,
	}
	return protoBp, nil
}

// UnmarshalOrchestration returns a orchestration from the orchestration protobuf encoding.
func UnmarshalOrchestration(p *pb.Orchestration) (
	*orchestration.Orchestration,
	error,
) {
	if ok, err := validate.ArgNotNil(p, "p"); !ok {
		return nil, err
	}

	links := make([]string, 0, len(p.Links))
	for _, l := range p.Links {
		links = append(links, l)
	}

	return &orchestration.Orchestration{
		Name:  p.Name,
		Links: links,
	}, nil
}
