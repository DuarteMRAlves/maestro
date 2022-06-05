package compiled

import (
	"github.com/DuarteMRAlves/maestro/internal"
)

// Pipeline defines an immutable pipeline that can be executed.
type Pipeline struct {
	name   internal.PipelineName
	mode   internal.ExecutionMode
	stages stageGraph
}

type stageGraph map[internal.StageName]*Stage

func (p *Pipeline) Name() internal.PipelineName {
	return p.name
}

func (p *Pipeline) Mode() internal.ExecutionMode {
	return p.mode
}

func (p *Pipeline) Stage(name internal.StageName) (*Stage, bool) {
	s, ok := p.stages[name]
	return s, ok
}

// Stage defines a step of a Pipeline
type Stage struct {
	name    internal.StageName
	method  internal.UnaryMethod
	inputs  []*internal.Link
	outputs []*internal.Link
}

func (s *Stage) Name() internal.StageName {
	return s.name
}

func (s *Stage) Method() internal.UnaryMethod {
	return s.method
}

func (s *Stage) Inputs() []*internal.Link {
	return s.inputs
}

func (s *Stage) Outputs() []*internal.Link {
	return s.outputs
}
