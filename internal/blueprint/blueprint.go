package blueprint

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/identifier"
)

const IdSize = 10

type Blueprint struct {
	Id     identifier.Id
	Name   string
	Stages []*Stage
	Links  []*Link
}

func (bp *Blueprint) Clone() *Blueprint {
	stages := make([]*Stage, 0, len(bp.Stages))
	for _, s := range bp.Stages {
		stages = append(stages, s.Clone())
	}

	links := make([]*Link, 0, len(bp.Links))
	for _, l := range bp.Links {
		links = append(links, l.Clone())
	}

	return &Blueprint{
		Id:     bp.Id.Clone(),
		Name:   bp.Name,
		Stages: stages,
		Links:  links,
	}
}

func (bp *Blueprint) String() string {
	return fmt.Sprintf(
		"Blueprint{Id:%v,Name:'%v',NumStages:%v,NumLinks:%v}",
		bp.Id,
		bp.Name,
		len(bp.Stages),
		len(bp.Links))
}
