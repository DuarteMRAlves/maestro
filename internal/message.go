package internal

import (
	"context"
)

type MessageField struct{ val string }

func (m MessageField) Unwrap() string { return m.val }

func (m MessageField) IsEmpty() bool { return m.val == "" }

func NewMessageField(field string) MessageField {
	return MessageField{val: field}
}

type OptionalMessageField struct {
	val     MessageField
	present bool
}

func (p OptionalMessageField) Unwrap() MessageField { return p.val }

func (p OptionalMessageField) Present() bool { return p.present }

func NewPresentMessageField(m MessageField) OptionalMessageField {
	return OptionalMessageField{val: m, present: true}
}

func NewEmptyMessageField() OptionalMessageField {
	return OptionalMessageField{}
}

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
