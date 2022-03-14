package internal

import (
	"regexp"
)

var linkNameRegExp, _ = regexp.Compile(`^[a-zA-Z0-9]+([-:_/][a-zA-Z0-9]+)*$|^$`)

type LinkName struct{ val string }

func (s LinkName) Unwrap() string { return s.val }

func (s LinkName) IsEmpty() bool { return s.val == "" }

func NewLinkName(name string) (LinkName, error) {
	if !isValidLinkName(name) {
		return LinkName{}, &InvalidIdentifier{Type: "link", Ident: name}
	}
	return LinkName{val: name}, nil
}

func isValidLinkName(name string) bool {
	return linkNameRegExp.MatchString(name)
}

type LinkEndpoint struct {
	stage StageName
	field OptionalMessageField
}

func (e LinkEndpoint) Stage() StageName {
	return e.stage
}

func (e LinkEndpoint) Field() OptionalMessageField {
	return e.field
}

func NewLinkEndpoint(
	stage StageName,
	field OptionalMessageField,
) LinkEndpoint {
	return LinkEndpoint{
		stage: stage,
		field: field,
	}
}

type Link struct {
	name          LinkName
	source        LinkEndpoint
	target        LinkEndpoint
	orchestration OrchestrationName
}

func (l Link) Name() LinkName {
	return l.name
}

func (l Link) Source() LinkEndpoint {
	return l.source
}

func (l Link) Target() LinkEndpoint {
	return l.target
}

func (l Link) Orchestration() OrchestrationName {
	return l.orchestration
}

func NewLink(
	name LinkName,
	source, target LinkEndpoint,
	orchestration OrchestrationName,
) Link {
	return Link{
		name:          name,
		source:        source,
		target:        target,
		orchestration: orchestration,
	}
}
