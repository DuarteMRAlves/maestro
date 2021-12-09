package blueprint

import (
	"fmt"
)

type Blueprint struct {
	Name  string
	Links []string
}

func (bp *Blueprint) Clone() *Blueprint {
	links := make([]string, 0, len(bp.Links))
	for _, l := range bp.Links {
		links = append(links, l)
	}

	return &Blueprint{
		Name:  bp.Name,
		Links: links,
	}
}

func (bp *Blueprint) String() string {
	return fmt.Sprintf(
		"Blueprint{Name:'%v',Links:%v}",
		bp.Name,
		bp.Links)
}
