package output

import (
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/flow/connection"
	"github.com/DuarteMRAlves/maestro/internal/flow/state"
)

// Output receives the output flow.State for a given stage and sends it to the
// next stages.
type Output interface {
	Yield(s *state.State)
}

// Cfg represents the several output connections for a stage
type Cfg struct {
	typ         Type
	connections []*connection.Connection
}

// Type defines the type of output the stage.Stage associated with this
// Cfg is expecting.
type Type string

const (
	// OutputInfer means the output type is not specified and should be inferred
	// from the received connections.
	OutputInfer Type = "Infer"
	// OutputSingle means the stage only sends its entire output to another
	// stage.
	OutputSingle Type = "Single"
	// OutputSink means the stage outputs the empty message and its output
	// should be discarded as it does not connect to any other stage.
	OutputSink Type = "Sink"
	// OutputSplit means the stage output should be split into multiple messages
	// sent to different stages.
	OutputSplit Type = "Split"
	// OutputDuplicate means the stage output should be duplicated and sent as
	// input to multiple stages.
	OutputDuplicate Type = "Duplicate"
)

func NewOutputCfg() *Cfg {
	return &Cfg{
		typ:         OutputInfer,
		connections: []*connection.Connection{},
	}
}

func (o *Cfg) Register(c *connection.Connection) error {
	for _, prev := range o.connections {
		if prev.HasSameLinkName(c) {
			return errdefs.InvalidArgumentWithMsg(
				"Link with an equal name already registered: %s",
				prev.LinkName())
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
