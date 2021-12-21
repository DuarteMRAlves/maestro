package reflection

import (
	"github.com/DuarteMRAlves/maestro/internal/validate"
	"github.com/jhump/protoreflect/desc"
)

// Message describes a grpc message
type Message interface {
	fullyQualifiedName() string
}

type message struct {
	desc *desc.MessageDescriptor
}

func newMessage(desc *desc.MessageDescriptor) (Message, error) {
	if ok, err := validate.ArgNotNil(desc, "desc"); !ok {
		return nil, err
	}
	s := &message{desc: desc}
	return s, nil
}

func (s *message) fullyQualifiedName() string {
	return s.desc.GetFullyQualifiedName()
}
