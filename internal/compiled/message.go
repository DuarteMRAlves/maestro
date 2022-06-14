package compiled

import (
	"context"
)

// MessageField specifies the name of a field in a message in the pipeline.
type MessageField string

// IsUnspecified reports whether this field is "".
func (m MessageField) IsUnspecified() bool { return m == "" }

// Message specifies an interface to send messages for the several stages.
type Message interface {
	SetField(MessageField, Message) error
	GetField(MessageField) (Message, error)
}

type EmptyMessageGen func() Message

type MessageDesc interface {
	Compatible(MessageDesc) bool
	EmptyGen() EmptyMessageGen
	GetField(MessageField) (MessageDesc, error)
}

type UnaryClientBuilder func(Address) (UnaryClient, error)

type UnaryClient interface {
	Call(ctx context.Context, req Message) (Message, error)
	Close() error
}

type UnaryMethod interface {
	ClientBuilder() UnaryClientBuilder
	Input() MessageDesc
	Output() MessageDesc
}
