package exec

import (
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/rpc"
)

// Input joins the input connections for a given stage and provides the next
// State to be processed.
type Input interface {
	Chan() <-chan *State
	Close()
	IsSource() bool
}

// SingleInput is a struct the implements the Input for a single input.
type SingleInput struct {
	conn *Link
}

func NewSingleInput(conn *Link) *SingleInput {
	i := &SingleInput{conn: conn}
	return i
}

func (i *SingleInput) Chan() <-chan *State {
	return i.conn.Chan()
}

func (i *SingleInput) Close() {}

func (i *SingleInput) IsSource() bool {
	return false
}

// SourceInput is the source of the orchestration. It defines the initial ids of
// the states and sends empty messages of the received type.
type SourceInput struct {
	id  int32
	msg rpc.Message
	ch  chan *State
	end chan struct{}
}

func NewSourceOutput(initial int32, msg rpc.Message) *SourceInput {
	ch := make(chan *State)

	i := &SourceInput{
		id:  initial,
		msg: msg,
		ch:  ch,
	}
	go func() {
		defer close(i.ch)
		defer close(i.end)
		for {
			select {
			case i.ch <- i.next():
			case <-i.end:
				return
			}
		}
	}()
	return i
}

func (i *SourceInput) next() *State {
	s := NewState(Id(i.id), i.msg.NewEmpty())
	i.id++
	return s
}

func (i *SourceInput) Chan() <-chan *State {
	return i.ch
}

func (i *SourceInput) Close() {
	i.end <- struct{}{}
}

func (i *SourceInput) IsSource() bool {
	return true
}

// InputBuilder registers the several connections for an input.
type InputBuilder struct {
	connections []*Link
	msg         rpc.Message
}

func NewInputBuilder() *InputBuilder {
	return &InputBuilder{
		connections: []*Link{},
	}
}

func (i *InputBuilder) WithMessage(msg rpc.Message) *InputBuilder {
	i.msg = msg
	return i
}

func (i *InputBuilder) WithConnection(c *Link) error {
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

func (i *InputBuilder) Build() (Input, error) {
	switch len(i.connections) {
	case 0:
		if i.msg == nil {
			return nil, errdefs.FailedPreconditionWithMsg(
				"message required without 0 connections",
			)
		}
		return NewSourceOutput(1, i.msg), nil
	case 1:
		return NewSingleInput(i.connections[0]), nil
	default:
		return nil, errdefs.FailedPreconditionWithMsg(
			"too many connections: expected 0 or 1 but received %d",
			len(i.connections),
		)
	}
}
