package orchestration

import (
	"fmt"
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
)

// Orchestration defines a set of links and stages that will be executed.
type Orchestration struct {
	// name defines a unique name for the orchestrations.
	name string
	// links specifies the names of the links contained in the orchestration.
	links []string
}

func New(name string, links []string) *Orchestration {
	return &Orchestration{
		name:  name,
		links: links,
	}
}

func (o *Orchestration) Name() string {
	return o.name
}

func (o *Orchestration) Links() []string {
	return o.links
}

func (o *Orchestration) Clone() *Orchestration {
	links := make([]string, 0, len(o.links))
	for _, l := range o.links {
		links = append(links, l)
	}

	return &Orchestration{
		name:  o.name,
		links: links,
	}
}

func (o *Orchestration) ToApi() *apitypes.Orchestration {
	links := make([]string, 0, len(o.links))
	for _, l := range o.links {
		links = append(links, l)
	}

	return &apitypes.Orchestration{
		Name:  o.name,
		Links: links,
	}
}

func (o *Orchestration) String() string {
	return fmt.Sprintf(
		"Orchestration{name:%v,links:%v}",
		o.name,
		o.links)
}
