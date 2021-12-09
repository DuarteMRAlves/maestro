package protobuff

import (
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/blueprint"
	"github.com/DuarteMRAlves/maestro/internal/validate"
)

func MarshalBlueprint(bp *blueprint.Blueprint) *pb.Blueprint {
	links := make([]string, 0, len(bp.Links))
	for _, l := range bp.Links {
		links = append(links, l)
	}
	return &pb.Blueprint{
		Name:  bp.Name,
		Links: links,
	}
}

func UnmarshalBlueprint(p *pb.Blueprint) (*blueprint.Blueprint, error) {
	if ok, err := validate.ArgNotNil(p, "p"); !ok {
		return nil, err
	}

	links, err := unmarshalLinks(p)
	if err != nil {
		return nil, err
	}

	return &blueprint.Blueprint{
		Name:  p.Name,
		Links: links,
	}, nil
}

func unmarshalLinks(p *pb.Blueprint) ([]string, error) {
	links := make([]string, 0, len(p.Links))
	for _, l := range p.Links {
		links = append(links, l)
	}
	return links, nil
}
