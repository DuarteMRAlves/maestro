package execution

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

func (c *Link) LinkName() api.LinkName {
	return c.link.Name
}

func (c *Link) HasSameLinkName(other *Link) bool {
	return c.link.Name == other.link.Name
}

func (c *Link) HasEmptyTargetField() bool {
	return c.link.TargetField == ""
}

func (c *Link) HasSameTargetField(other *Link) bool {
	return c.link.TargetField == other.link.TargetField
}

func (c *Link) Chan() chan *State {
	return c.ch
}
