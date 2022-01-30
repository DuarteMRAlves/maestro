package execution

import (
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
)

// Input joins the input connections for a given stage and provides the next
// State to be processed.
type Input interface {
	Next() (*State, error)
}

// InputCfg represents the several input connections for a stage
type InputCfg struct {
	connections []*Connection
}

func NewInputCfg() *InputCfg {
	return &InputCfg{
		connections: []*Connection{},
	}
}

func (i *InputCfg) Register(c *Connection) error {
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

func (i *InputCfg) UnregisterIfExists(search *Connection) {
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

func (i *InputCfg) ToInput() Input {
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
