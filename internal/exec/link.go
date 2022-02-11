package exec

import (
	"github.com/DuarteMRAlves/maestro/internal/api"
)

// Link is a connection between two stages where data is transferred.
type Link struct {
	link *api.Link
	ch   chan *State
}

func NewLink(l *api.Link) *Link {
	f := &Link{
		link: l,
		ch:   make(chan *State),
	}
	return f
}

func (l *Link) LinkName() api.LinkName {
	return l.link.Name
}

func (l *Link) HasSameLinkName(other *Link) bool {
	return l.link.Name == other.link.Name
}

func (l *Link) SourceField() string {
	return l.link.SourceField
}
func (l *Link) TargetField() string {
	return l.link.TargetField
}

func (l *Link) HasEmptyTargetField() bool {
	return l.link.TargetField == ""
}

func (l *Link) HasSameTargetField(other *Link) bool {
	return l.link.TargetField == other.link.TargetField
}

func (l *Link) Chan() chan *State {
	return l.ch
}
