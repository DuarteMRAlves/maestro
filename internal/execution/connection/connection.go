package connection

import (
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"github.com/DuarteMRAlves/maestro/internal/execution/state"
	"github.com/DuarteMRAlves/maestro/internal/queue"
	"github.com/DuarteMRAlves/maestro/internal/storage"
)

// Connection is a connection between two stages where data is transferred.
type Connection struct {
	link  *storage.Link
	queue queue.Ring
}

func New(l *storage.Link) (*Connection, error) {
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

func (c *Connection) LinkName() apitypes.LinkName {
	return c.link.Name()
}

func (c *Connection) HasSameLinkName(other *Connection) bool {
	return c.link.Name() == other.link.Name()
}

func (c *Connection) HasEmptyTargetField() bool {
	return c.link.TargetField() == ""
}

func (c *Connection) HasSameTargetField(other *Connection) bool {
	return c.link.TargetField() == other.link.TargetField()
}

func (c *Connection) Push(s *state.State) {
	c.queue.Push(s)
}

func (c *Connection) Pop() *state.State {
	return c.queue.Pop().(*state.State)
}
