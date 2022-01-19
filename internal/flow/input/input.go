package input

import (
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/flow/connection"
	"github.com/DuarteMRAlves/maestro/internal/flow/state"
)

// Input joins the input connections for a given stage and provides the next
// State to be processed.
type Input interface {
	Next() (*state.State, error)
}

// Cfg represents the several input connections for a stage
type Cfg struct {
	typ         Type
	connections []*connection.Connection
}

// Type defines the type of input that the stage.Stage associated with this
// Cfg is expecting.
type Type string

const (
	// InputInfer means the input type is not specified and should be inferred
	// from the received connections.
	InputInfer Type = "Infer"
	// InputSingle means the stage only receives input from another stage.
	InputSingle Type = "Single"
	// InputSource means the stage is a source stage and receives as input an
	// empty message that should be generated by the orchestrator.
	InputSource Type = "Source"
	// InputMerge means the input is a merge of multiple different messages,
	// coming from different stages.
	InputMerge Type = "Merge"
	// InputCollect means the input collects from multiple stages, but the
	// received messages are all equal and should be directly sent to this
	// input stage.
	InputCollect Type = "Collect"
)

func NewInputCfg() *Cfg {
	return &Cfg{
		typ:         InputInfer,
		connections: []*connection.Connection{},
	}
}

func (i *Cfg) Register(c *connection.Connection) error {
	// A previous link that consumes the entire message already exists
	if len(i.connections) == 1 && i.connections[0].HasEmptyTargetField() {
		return errdefs.FailedPreconditionWithMsg(
			"link that receives the full message already exists")
	}
	for _, prev := range i.connections {
		if prev.HasSameTargetField(c) {
			return errdefs.InvalidArgumentWithMsg(
				"link with the same target field already registered: %s",
				prev.LinkName())
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
