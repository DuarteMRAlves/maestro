package execution

import (
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
)

// Output receives the output flow.State for a given stage and sends it to the
// next stages.
type Output interface {
	Yield(s *State)
	IsSink() bool
}

// SingleOutput is a struct that implements Output for a single output
// connection.
type SingleOutput struct {
	connection *Connection
}

func (o *SingleOutput) Yield(s *State) {
	o.connection.Push(s)
}

func (o *SingleOutput) IsSink() bool {
	return false
}

// SinkOutput defines the last output of the orchestration, where all messages
// are dropped.
type SinkOutput struct {
}

func (o *SinkOutput) Yield(_ *State) {
	// Do nothing, just drop message.
}

func (o *SinkOutput) IsSink() bool {
	return true
}

// OutputBuilder registers the several connections for an output.
type OutputBuilder struct {
	connections []*Connection
}

func NewOutputBuilder() *OutputBuilder {
	return &OutputBuilder{
		connections: []*Connection{},
	}
}

func (o *OutputBuilder) WithConnection(c *Connection) error {
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

func (o *OutputBuilder) UnregisterIfExists(search *Connection) {
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

func (o *OutputBuilder) Build() (Output, error) {
	switch len(o.connections) {
	case 0:
		return &SinkOutput{}, nil
	case 1:
		return &SingleOutput{connection: o.connections[0]}, nil
	default:
		return nil, errdefs.FailedPreconditionWithMsg(
			"too many connections: expected 0 or 1 but received %d",
			len(o.connections),
		)
	}
}
