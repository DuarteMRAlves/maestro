package orchestration

import (
	"fmt"
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
)

type Orchestration struct {
	Name  string
	Links []string
}

func (o *Orchestration) Clone() *Orchestration {
	links := make([]string, 0, len(o.Links))
	for _, l := range o.Links {
		links = append(links, l)
	}

	return &Orchestration{
		Name:  o.Name,
		Links: links,
	}
}

func (o *Orchestration) ToApi() *apitypes.Orchestration {
	links := make([]string, 0, len(o.Links))
	for _, l := range o.Links {
		links = append(links, l)
	}

	return &apitypes.Orchestration{
		Name:  o.Name,
		Links: links,
	}
}

func (o *Orchestration) String() string {
	return fmt.Sprintf(
		"Orchestration{Name:%v,Links:%v}",
		o.Name,
		o.Links)
}
