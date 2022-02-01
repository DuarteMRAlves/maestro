package execution

import (
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/rpc"
	"sync/atomic"
)

// Input joins the input connections for a given stage and provides the next
// State to be processed.
type Input interface {
	Next() (*State, error)
	IsSource() bool
}

// InputBuilder registers the several connections for an input.
type InputBuilder struct {
	connections []*Connection
	msg         rpc.Message
}

func NewInputBuilder() *InputBuilder {
	return &InputBuilder{
		connections: []*Connection{},
	}
}

func (i *InputBuilder) WithMessage(msg rpc.Message) *InputBuilder {
	i.msg = msg
	return i
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
	case 0:
		return &SourceInput{id: 0, msg: i.msg}
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

func (i *SingleInput) IsSource() bool {
	return false
}

// SourceInput is the source of the orchestration. It defines the initial ids of
// the states and sends empty messages of the received type.
type SourceInput struct {
	id  int32
	msg rpc.Message
}

func (i *SourceInput) Next() (*State, error) {
	id := Id(atomic.AddInt32(&(i.id), 1))
	return NewState(id, i.msg.NewEmpty()), nil
}

func (i *SourceInput) IsSource() bool {
	return true
}
