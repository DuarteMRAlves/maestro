package storage

import (
	"fmt"
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
)

// Orchestration defines a set of links and stages that will be executed.
type Orchestration struct {
	// name defines a unique name for the orchestrations.
	name apitypes.OrchestrationName
	// phase defines the current phase of this Orchestration.
	phase apitypes.OrchestrationPhase
	// stages specifies the stages that will be executed by this Orchestration.
	stages []*Stage
	// links specifies the links contained in this Orchestration.
	links []*Link
}

func NewOrchestration(name apitypes.OrchestrationName) *Orchestration {
	return &Orchestration{
		name:   name,
		phase:  apitypes.OrchestrationPending,
		stages: []*Stage{},
		links:  []*Link{},
	}
}

func (o *Orchestration) Name() apitypes.OrchestrationName {
	return o.name
}

func (o *Orchestration) Links() []*Link {
	return o.links
}

func (o *Orchestration) Clone() *Orchestration {
	links := make([]*Link, 0, len(o.links))
	for _, l := range o.links {
		links = append(links, l)
	}

	return &Orchestration{
		name:  o.name,
		phase: o.phase,
		links: links,
	}
}

func (o *Orchestration) AddStage(s *Stage) {
	o.stages = append(o.stages, s)
}

func (o *Orchestration) AddLink(l *Link) {
	o.links = append(o.links, l)
}

func (o *Orchestration) ToApi() *apitypes.Orchestration {
	links := make([]apitypes.LinkName, 0, len(o.links))
	for _, l := range o.links {
		links = append(links, l.Name())
	}

	return &apitypes.Orchestration{
		Name:  o.name,
		Phase: o.phase,
		Links: links,
	}
}

func (o *Orchestration) String() string {
	return fmt.Sprintf(
		"Orchestration{name:%v,phase:%v,links:%v}",
		o.name,
		o.phase,
		o.links,
	)
}
