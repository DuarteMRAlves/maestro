package execution

import (
	"github.com/DuarteMRAlves/maestro/internal/api"
)

// Connection is a connection between two stages where data is transferred.
type Connection struct {
	link *api.Link
	ch   chan *State
}

func NewConnection(l *api.Link) *Connection {
	f := &Connection{
		link: l,
		ch:   make(chan *State),
	}
	return f
}

func (c *Connection) LinkName() api.LinkName {
	return c.link.Name
}

func (c *Connection) HasSameLinkName(other *Connection) bool {
	return c.link.Name == other.link.Name
}

func (c *Connection) HasEmptyTargetField() bool {
	return c.link.TargetField == ""
}

func (c *Connection) HasSameTargetField(other *Connection) bool {
	return c.link.TargetField == other.link.TargetField
}

func (c *Connection) Chan() chan *State {
	return c.ch
}
