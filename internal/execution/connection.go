package execution

import (
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/queue"
)

// Connection is a connection between two stages where data is transferred.
type Connection struct {
	link  *api.Link
	queue queue.Ring
}

func NewConnection(l *api.Link) (*Connection, error) {
	q, err := queue.NewRing(1)
	if err != nil {
		return nil, err
	}
	f := &Connection{
		link:  l,
		queue: q,
	}
	return f, nil
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

func (c *Connection) Push(s *State) {
	c.queue.Push(s)
}

func (c *Connection) Pop() *State {
	return c.queue.Pop().(*State)
}
