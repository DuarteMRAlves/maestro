package exec

import (
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
)

// Output receives the output flow.State for a given stage and sends it to the
// next stages.
type Output interface {
	Chan() chan<- *State
	Close()
	IsSink() bool
}

// SingleOutput is a struct that implements Output for a single output
// connection.
type SingleOutput struct {
	connection *Link
}

func NewSingleOutput(conn *Link) *SingleOutput {
	o := &SingleOutput{
		connection: conn,
	}
	return o
}

func (o *SingleOutput) Chan() chan<- *State {
	return o.connection.Chan()
}

func (o *SingleOutput) Close() {
}

func (o *SingleOutput) IsSink() bool {
	return false
}

// SinkOutput defines the last output of the orchestration, where all messages
// are dropped.
type SinkOutput struct {
	ch  chan *State
	end chan struct{}
}

func NewSinkOutput() *SinkOutput {
	ch := make(chan *State)
	end := make(chan struct{})

	o := &SinkOutput{
		ch:  ch,
		end: end,
	}
	go func() {
		defer close(o.ch)
		defer close(o.end)
		for {
			select {
			// Discard results
			case <-o.ch:
			case <-o.end:
				return
			}
		}
	}()
	return o
}

func (o *SinkOutput) Chan() chan<- *State {
	return o.ch
}

func (o *SinkOutput) Close() {
	o.end <- struct{}{}
}

func (o *SinkOutput) IsSink() bool {
	return true
}

// OutputBuilder registers the several connections for an output.
type OutputBuilder struct {
	connections []*Link
}

func NewOutputBuilder() *OutputBuilder {
	return &OutputBuilder{
		connections: []*Link{},
	}
}

func (o *OutputBuilder) WithConnection(c *Link) error {
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

func (o *OutputBuilder) UnregisterIfExists(search *Link) {
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
		return NewSinkOutput(), nil
	case 1:
		return NewSingleOutput(o.connections[0]), nil
	default:
		return nil, errdefs.FailedPreconditionWithMsg(
			"too many connections: expected 0 or 1 but received %d",
			len(o.connections),
		)
	}
}
