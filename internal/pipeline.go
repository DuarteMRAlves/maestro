package internal

import (
	"fmt"
	"regexp"
)

var pipelineNameReqExp, _ = regexp.Compile(`^[a-zA-Z0-9]+([-:_/][a-zA-Z0-9]+)*$|^$`)

type PipelineName struct{ val string }

var emptyPipelineName = PipelineName{val: ""}

func (o PipelineName) Unwrap() string { return o.val }

func (o PipelineName) IsEmpty() bool { return o.val == emptyPipelineName.val }

func (o PipelineName) String() string {
	return o.val
}

func NewPipelineName(name string) (PipelineName, error) {
	if isValidPipelineName(name) {
		return PipelineName{name}, nil
	}
	return emptyPipelineName, &invalidPipelineName{name: name}
}

func isValidPipelineName(name string) bool {
	return pipelineNameReqExp.MatchString(name)
}

type invalidPipelineName struct{ name string }

func (err *invalidPipelineName) Error() string {
	return fmt.Sprintf("invalid pipeline name: '%s'", err.name)
}

type ExecutionMode int

const (
	// OfflineExecution is the default execution mode.
	OfflineExecution ExecutionMode = iota
	OnlineExecution
)

func (e ExecutionMode) String() string {
	switch e {
	case OfflineExecution:
		return "OfflineExecution"
	case OnlineExecution:
		return "OnlineExecution"
	default:
		return "UnknownExecutionMode"
	}
}

type Pipeline struct {
	name   PipelineName
	stages []StageName
	links  []LinkName
	mode   ExecutionMode
}

func (o Pipeline) Name() PipelineName {
	return o.name
}

func (o Pipeline) Stages() []StageName {
	return o.stages
}

func (o Pipeline) Links() []LinkName {
	return o.links
}

func (o Pipeline) Mode() ExecutionMode {
	return o.mode
}

type PipelineOpt func(*Pipeline)

func WithStages(ss ...StageName) PipelineOpt {
	return func(p *Pipeline) {
		p.stages = ss
	}
}

func WithLinks(ll ...LinkName) PipelineOpt {
	return func(p *Pipeline) {
		p.links = ll
	}
}

func WithOfflineExec() PipelineOpt {
	return func(p *Pipeline) {
		p.mode = OfflineExecution
	}
}

func WithOnlineExec() PipelineOpt {
	return func(p *Pipeline) {
		p.mode = OnlineExecution
	}
}

func NewPipeline(name PipelineName, opts ...PipelineOpt) Pipeline {
	p := Pipeline{name: name}
	for _, o := range opts {
		o(&p)
	}
	return p
}

func FromPipeline(src Pipeline, opts ...PipelineOpt) Pipeline {
	p := Pipeline{
		name:   src.name,
		stages: src.stages,
		links:  src.links,
		mode:   src.mode,
	}
	for _, o := range opts {
		o(&p)
	}
	return p
}
