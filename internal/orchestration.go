package internal

import (
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"regexp"
)

var orchestrationNameReqExp, _ = regexp.Compile(`^[a-zA-Z0-9]+([-:_/][a-zA-Z0-9]+)*$|^$`)

type OrchestrationName struct{ val string }

var emptyOrchestrationName = OrchestrationName{val: ""}

func (o OrchestrationName) Unwrap() string { return o.val }

func (o OrchestrationName) IsEmpty() bool { return o.val == emptyOrchestrationName.val }

func NewOrchestrationName(name string) (OrchestrationName, error) {
	if isValidOrchestrationName(name) {
		return OrchestrationName{name}, nil
	}
	err := errdefs.InvalidArgumentWithMsg("invalid name '%v'", name)
	return emptyOrchestrationName, err
}

func isValidOrchestrationName(name string) bool {
	return orchestrationNameReqExp.MatchString(name)
}

type Orchestration struct {
	name   OrchestrationName
	stages []StageName
	links  []LinkName
}

func (o Orchestration) Name() OrchestrationName {
	return o.name
}

func (o Orchestration) Stages() []StageName {
	return o.stages
}

func (o Orchestration) Links() []LinkName {
	return o.links
}

func NewOrchestration(
	name OrchestrationName,
	stages []StageName,
	links []LinkName,
) Orchestration {
	return Orchestration{
		name:   name,
		stages: stages,
		links:  links,
	}
}
