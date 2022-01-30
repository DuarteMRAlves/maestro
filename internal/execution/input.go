package execution

import (
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
)

// Input joins the input connections for a given stage and provides the next
// State to be processed.
type Input interface {
	Next() (*State, error)
}

// InputBuilder registers the several connections for an input.
type InputBuilder struct {
	connections []*Connection
}

func NewInputBuilder() *InputBuilder {
	return &InputBuilder{
		connections: []*Connection{},
	}
}

func (i *InputBuilder) WithConnection(c *Connection) error {
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

func (i *InputBuilder) Build() Input {
	switch len(i.connections) {
	case 1:
		return &SingleInput{connection: i.connections[0]}
	}
	return nil
}

// SingleInput is a struct the implements the Input for a single input.
type SingleInput struct {
	connection *Connection
}

func (i *SingleInput) Next() (*State, error) {
	return i.connection.Pop(), nil
}
