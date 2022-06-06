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

// StageVisitor is a function to process stages.
type StageVisitor func(s *Stage) error

// LinkVisitor is a function to process links.
type LinkVisitor func(l *internal.Link) error

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

// VisitStages iterates through the stages in the pipeline and executes the
// visitor function. Every stage is only seen once. If an error is returned
// by the visitor function, the iteration is stopped and the error is returned.
func (p *Pipeline) VisitStages(v StageVisitor) error {
	for _, s := range p.stages {
		if err := v(s); err != nil {
			return err
		}
	}
	return nil
}

// VisitLinks iterates through the links in the pipeline and executes the
// visitor function. Every link is only seen once. If an error is returned
// by the visitor function, the iteration is stopped and the error is returned.
func (p *Pipeline) VisitLinks(v LinkVisitor) error {
	for _, s := range p.stages {
		for _, l := range s.inputs {
			if err := v(l); err != nil {
				return err
			}
		}
	}
	return nil
}

// Stage defines a step of a Pipeline
type Stage struct {
	name    internal.StageName
	address internal.Address
	method  internal.UnaryMethod
	inputs  []*internal.Link
	outputs []*internal.Link
}

func (s *Stage) Name() internal.StageName {
	return s.name
}

func (s *Stage) Address() internal.Address {
	return s.address
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
