package compiled

import (
	"fmt"

	"github.com/DuarteMRAlves/maestro/internal/message"
)

const defaultLinkSize uint = 10

type Link struct {
	name             LinkName
	source           *LinkEndpoint
	target           *LinkEndpoint
	numEmptyMessages uint
	size             uint
}

func (l *Link) Name() LinkName {
	if l == nil {
		return LinkName{}
	}
	return l.name
}

func (l *Link) Source() *LinkEndpoint {
	if l == nil {
		return nil
	}
	return l.source
}

func (l *Link) Target() *LinkEndpoint {
	if l == nil {
		return nil
	}
	return l.target
}

func (l *Link) Size() uint {
	if l == nil {
		return defaultLinkSize
	}
	return l.size
}

func (l *Link) NumEmptyMessages() uint {
	if l == nil {
		return 0
	}
	return l.numEmptyMessages
}

func NewLink(
	name LinkName,
	source, target *LinkEndpoint,
	size uint,
	numEmptyMessages uint,
) *Link {
	fmt.Println("Creating link with size", size)
	return &Link{
		name:             name,
		source:           source,
		target:           target,
		size:             size,
		numEmptyMessages: numEmptyMessages,
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
	field message.Field
}

func (e *LinkEndpoint) Stage() StageName {
	if e == nil {
		return StageName{}
	}
	return e.stage
}

func (e *LinkEndpoint) Field() message.Field {
	if e == nil {
		return ""
	}
	return e.field
}

func NewLinkEndpoint(stage StageName, field message.Field) LinkEndpoint {
	return LinkEndpoint{
		stage: stage,
		field: field,
	}
}
