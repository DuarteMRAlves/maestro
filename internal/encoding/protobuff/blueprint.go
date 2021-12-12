package protobuff

import (
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/blueprint"
	"github.com/DuarteMRAlves/maestro/internal/validate"
)

// MarshalBlueprint returns a protobuf encoding of the given blueprint.
func MarshalBlueprint(bp *blueprint.Blueprint) (*pb.Blueprint, error) {
	if ok, err := validate.ArgNotNil(bp, "bp"); !ok {
		return nil, err
	}
	links := make([]string, 0, len(bp.Links))
	for _, l := range bp.Links {
		links = append(links, l)
	}
	protoBp := &pb.Blueprint{
		Name:  bp.Name,
		Links: links,
	}
	return protoBp, nil
}

// UnmarshalBlueprint returns a blueprint from the blueprint protobuf encoding.
func UnmarshalBlueprint(p *pb.Blueprint) (*blueprint.Blueprint, error) {
	if ok, err := validate.ArgNotNil(p, "p"); !ok {
		return nil, err
	}

	links := make([]string, 0, len(p.Links))
	for _, l := range p.Links {
		links = append(links, l)
	}

	return &blueprint.Blueprint{
		Name:  p.Name,
		Links: links,
	}, nil
}
