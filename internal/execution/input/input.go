package input

import (
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/execution/connection"
	"github.com/DuarteMRAlves/maestro/internal/execution/state"
)

// Input joins the input connections for a given stage and provides the next
// State to be processed.
type Input interface {
	Next() (*state.State, error)
}

// Cfg represents the several input connections for a stage
type Cfg struct {
	connections []*connection.Connection
}

func NewCfg() *Cfg {
	return &Cfg{
		connections: []*connection.Connection{},
	}
}

func (i *Cfg) Register(c *connection.Connection) error {
	// A previous link that consumes the entire message already exists
	if len(i.connections) == 1 && i.connections[0].HasEmptyTargetField() {
		return errdefs.FailedPreconditionWithMsg(
			"link that receives the full message already exists",
		)
	}
	for _, prev := range i.connections {
		if prev.HasSameTargetField(c) {
			return errdefs.InvalidArgumentWithMsg(
				"link with the same target field already registered: %s",
				prev.LinkName(),
			)
		}
	}
	i.connections = append(i.connections, c)
	return nil
}

func (i *Cfg) UnregisterIfExists(search *connection.Connection) {
	idx := -1
	for j, c := range i.connections {
		if c.HasSameLinkName(search) {
			idx = j
			break
		}
	}
	if idx != -1 {
		i.connections[idx] = i.connections[len(i.connections)-1]
		i.connections = i.connections[:len(i.connections)-1]
	}
}

func (i *Cfg) ToInput() Input {
	switch len(i.connections) {
	case 1:
		return &SingleInput{connection: i.connections[0]}
	}
	return nil
}

// SingleInput is a struct the implements the Input for a single input.
type SingleInput struct {
	connection *connection.Connection
}

func (i *SingleInput) Next() (*state.State, error) {
	return i.connection.Pop(), nil
}
