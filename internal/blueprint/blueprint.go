package blueprint

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/identifier"
)

const IdSize = 10

type Blueprint struct {
	Id     identifier.Id
	Name   string
	stages []*Stage
	links  []*Link
}

func (bp *Blueprint) Clone() *Blueprint {
	stages := make([]*Stage, 0, len(bp.stages))
	for _, s := range bp.stages {
		stages = append(stages, s.Clone())
	}

	links := make([]*Link, 0, len(bp.links))
	for _, l := range bp.links {
		links = append(links, l.Clone())
	}

	return &Blueprint{
		Id:     bp.Id.Clone(),
		Name:   bp.Name,
		stages: stages,
		links:  links,
	}
}

func (bp *Blueprint) String() string {
	return fmt.Sprintf(
		"Blueprint{Id:%v,Name:'%v',NumStages:%v,NumLinks:%v}",
		bp.Id,
		bp.Name,
		len(bp.stages),
		len(bp.links))
}
