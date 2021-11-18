package protobuff

import (
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/assert"
	"github.com/DuarteMRAlves/maestro/internal/blueprint"
)

func MarshalBlueprint(bp *blueprint.Blueprint) *pb.Blueprint {
	stages := make([]*pb.Stage, 0, len(bp.Stages))
	for _, s := range bp.Stages {
		stages = append(stages, MarshalStage(s))
	}
	links := make([]*pb.Link, 0, len(bp.Links))
	for _, l := range bp.Links {
		links = append(links, MarshalLink(l))
	}
	return &pb.Blueprint{
		Name:   bp.Name,
		Stages: stages,
		Links:  links,
	}
}

func UnmarshalBlueprint(p *pb.Blueprint) (*blueprint.Blueprint, error) {
	if ok, err := assert.ArgNotNil(p, "p"); !ok {
		return nil, err
	}

	stages, err := unmarshalStages(p)
	if err != nil {
		return nil, err
	}

	links, err := unmarshalLinks(p)
	if err != nil {
		return nil, err
	}

	return &blueprint.Blueprint{
		Name:   p.Name,
		Stages: stages,
		Links:  links,
	}, nil
}

func unmarshalStages(p *pb.Blueprint) ([]*blueprint.Stage, error) {
	stages := make([]*blueprint.Stage, 0, len(p.Stages))
	for _, s := range p.Stages {
		res, err := UnmarshalStage(s)
		if err != nil {
			return nil, err
		}
		stages = append(stages, res)
	}
	return stages, nil
}

func unmarshalLinks(p *pb.Blueprint) ([]*blueprint.Link, error) {
	links := make([]*blueprint.Link, 0, len(p.Links))
	for _, l := range p.Links {
		res, err := UnmarshalLink(l)
		if err != nil {
			return nil, err
		}
		links = append(links, res)
	}
	return links, nil
}
