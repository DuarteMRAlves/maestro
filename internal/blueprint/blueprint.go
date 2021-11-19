package blueprint

import (
	"fmt"
)

const IdSize = 10

type Blueprint struct {
	Name   string
	Stages []string
	Links  []string
}

func (bp *Blueprint) Clone() *Blueprint {
	stages := make([]string, 0, len(bp.Stages))
	for _, s := range bp.Stages {
		stages = append(stages, s)
	}

	links := make([]string, 0, len(bp.Links))
	for _, l := range bp.Links {
		links = append(links, l)
	}

	return &Blueprint{
		Name:   bp.Name,
		Stages: stages,
		Links:  links,
	}
}

func (bp *Blueprint) String() string {
	return fmt.Sprintf(
		"Blueprint{Name:'%v',Stages:%v,Links:%v}",
		bp.Name,
		bp.Stages,
		bp.Links)
}
