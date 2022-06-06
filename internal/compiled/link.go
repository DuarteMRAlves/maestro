package compiled

import (
	"fmt"
)

type Link struct {
	name   LinkName
	source LinkEndpoint
	target LinkEndpoint
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

func NewLink(
	name LinkName,
	source, target LinkEndpoint,
) Link {
	return Link{
		name:   name,
		source: source,
		target: target,
	}
}

type LinkName struct{ val string }

func (l LinkName) Unwrap() string { return l.val }

func (l LinkName) IsEmpty() bool { return l.val == "" }

func (l LinkName) String() string {
	return l.val
}

func NewLinkName(name string) (LinkName, error) {
	if !validateResourceName(name) {
		return LinkName{}, &invalidLinkName{name: name}
	}
	return LinkName{val: name}, nil
}

type invalidLinkName struct{ name string }

func (err *invalidLinkName) Error() string {
	return fmt.Sprintf("invalid link name: '%s'", err.name)
}

type LinkEndpoint struct {
	stage StageName
	field MessageField
}

func (e LinkEndpoint) Stage() StageName {
	return e.stage
}

func (e LinkEndpoint) Field() MessageField {
	return e.field
}

func NewLinkEndpoint(stage StageName, field MessageField) LinkEndpoint {
	return LinkEndpoint{
		stage: stage,
		field: field,
	}
}