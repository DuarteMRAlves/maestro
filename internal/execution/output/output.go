package output

import (
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/execution/connection"
	"github.com/DuarteMRAlves/maestro/internal/execution/state"
)

// Output receives the output flow.State for a given stage and sends it to the
// next stages.
type Output interface {
	Yield(s *state.State)
}

// Cfg represents the several output connections for a stage
type Cfg struct {
	connections []*connection.Connection
}

func NewCfg() *Cfg {
	return &Cfg{
		connections: []*connection.Connection{},
	}
}

func (o *Cfg) Register(c *connection.Connection) error {
	for _, prev := range o.connections {
		if prev.HasSameLinkName(c) {
			return errdefs.InvalidArgumentWithMsg(
				"Link with an equal name already registered: %s",
				prev.LinkName(),
			)
		}
	}

	o.connections = append(o.connections, c)
	return nil
}

func (o *Cfg) UnregisterIfExists(search *connection.Connection) {
	idx := -1
	for i, c := range o.connections {
		if c.HasSameLinkName(search) {
			idx = i
			break
		}
	}
	if idx != -1 {
		o.connections[idx] = o.connections[len(o.connections)-1]
		o.connections = o.connections[:len(o.connections)-1]
	}
}

func (o *Cfg) ToFlow() Output {
	switch len(o.connections) {
	case 1:
		return &SingleOutput{connection: o.connections[0]}
	}
	return nil
}

// SingleOutput is a struct that implements Output for a single output
// connection.
type SingleOutput struct {
	connection *connection.Connection
}

func (o *SingleOutput) Yield(s *state.State) {
	o.connection.Push(s)
}
