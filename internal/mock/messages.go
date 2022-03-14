package mock

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal"
)

type Message struct {
	Fields map[internal.MessageField]interface{}
}

func (m *Message) SetField(
	field internal.MessageField,
	val internal.Message,
) error {
	m.Fields[field] = val
	return nil
}

func (m *Message) GetField(field internal.MessageField) (
	internal.Message,
	error,
) {
	f, ok := m.Fields[field]
	if !ok {
		return nil, fmt.Errorf("field '%s' does not exist", field)
	}
	msg, ok := f.(*Message)
	if !ok {
		return nil, fmt.Errorf("field '%s' is not a message", field)
	}
	return msg, nil
}

func NewGen() internal.EmptyMessageGen {
	return func() internal.Message {
		return &Message{Fields: map[internal.MessageField]interface{}{}}
	}
}
