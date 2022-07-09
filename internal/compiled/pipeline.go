package compiled

import (
	"fmt"

	"github.com/DuarteMRAlves/maestro/internal/api"
)

// Pipeline defines an immutable pipeline that can be executed.
type Pipeline struct {
	name   PipelineName
	mode   ExecutionMode
	stages stageGraph
}

// StageVisitor is a function to process stages.
type StageVisitor func(s *Stage) error

// LinkVisitor is a function to process links.
type LinkVisitor func(l *Link) error

type stageGraph map[StageName]*Stage

func (p *Pipeline) Name() PipelineName {
	return p.name
}

func (p *Pipeline) Mode() ExecutionMode {
	return p.mode
}

func (p *Pipeline) Stage(name StageName) (*Stage, bool) {
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

type PipelineName struct{ val string }

var emptyPipelineName = PipelineName{val: ""}

func (o PipelineName) Unwrap() string { return o.val }

func (o PipelineName) IsEmpty() bool { return o.val == emptyPipelineName.val }

func (o PipelineName) String() string {
	return o.val
}

func NewPipelineName(name string) (PipelineName, error) {
	if validateResourceName(name) {
		return PipelineName{name}, nil
	}
	return emptyPipelineName, &invalidPipelineName{name: name}
}

type invalidPipelineName struct{ name string }

func (err *invalidPipelineName) Error() string {
	return fmt.Sprintf("invalid pipeline name: '%s'", err.name)
}

type ExecutionMode struct{ val int }

const (
	// OfflineExecution is the default execution mode.
	offlineExecution int = iota
	onlineExecution
)

var (
	OfflineExecution = ExecutionMode{val: offlineExecution}
	OnlineExecution  = ExecutionMode{val: onlineExecution}
)

func modeFromApi(m api.ExecutionMode) (ExecutionMode, error) {
	switch m {
	case api.OfflineExecution:
		return OfflineExecution, nil
	case api.OnlineExecution:
		return OnlineExecution, nil
	default:
		return OfflineExecution, fmt.Errorf("unknown execution mode: %v", m)
	}
}

func (e ExecutionMode) String() string {
	switch e.val {
	case offlineExecution:
		return "OfflineExecution"
	case onlineExecution:
		return "OnlineExecution"
	default:
		return "UnknownExecutionMode"
	}
}
