package domain

import (
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"regexp"
)

var orchestrationNameReqExp, _ = regexp.Compile(`^[a-zA-Z0-9]+([-:_/][a-zA-Z0-9]+)*$`)

type orchestrationName string

func (o orchestrationName) Unwrap() string {
	return string(o)
}

func NewOrchestrationName(name string) (OrchestrationName, error) {
	if isValidOrchestrationName(name) {
		return orchestrationName(name), nil
	}
	return nil, errdefs.InvalidArgumentWithMsg("invalid name '%v'", name)
}

func isValidOrchestrationName(name string) bool {
	return orchestrationNameReqExp.MatchString(name)
}

type orchestration struct {
	name   OrchestrationName
	stages []Stage
	links  []Link
}

func (o orchestration) Name() OrchestrationName {
	return o.name
}

func (o orchestration) Stages() []Stage {
	return o.stages
}

func (o orchestration) Links() []Link {
	return o.links
}

func NewOrchestration(name OrchestrationName) Orchestration {
	return &orchestration{
		name:   name,
		stages: []Stage{},
		links:  []Link{},
	}
}
