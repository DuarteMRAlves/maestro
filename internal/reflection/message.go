package reflection

import (
	"github.com/DuarteMRAlves/maestro/internal/validate"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jhump/protoreflect/desc"
)

// Message describes a grpc message
type Message interface {
	FullyQualifiedName() string
	// Compatible verifies if two Messages are compatible, meaning fields with
	// equal numbers have the same type.
	Compatible(other Message) bool
	// GetMessageField searches for a message field with the given name. It
	// returns a Message with the given field and true as the second value. If
	// no field with the given name exists, or the field has the wrong type, nil
	// and false is returned.
	GetMessageField(name string) (Message, bool)
}

type message struct {
	desc *desc.MessageDescriptor
}

func NewMessage(desc *desc.MessageDescriptor) (Message, error) {
	if ok, err := validate.ArgNotNil(desc, "desc"); !ok {
		return nil, err
	}
	s := newMessageInternal(desc)
	return s, nil
}

func newMessageInternal(desc *desc.MessageDescriptor) Message {
	s := &message{desc: desc}
	return s
}

func (s *message) FullyQualifiedName() string {
	return s.desc.GetFullyQualifiedName()
}

func (s *message) Compatible(other Message) bool {
	otherMsg, ok := other.(*message)
	if !ok {
		return false
	}
	return s.cmpFields(otherMsg)
}

func (s *message) cmpFields(o *message) bool {
	for _, sField := range s.desc.GetFields() {
		number := sField.GetNumber()
		oField := o.desc.FindFieldByNumber(number)

		// Ignore unmatched fields
		if oField == nil {
			continue
		}

		// Both fields must be repeatable or not repeatable
		if sField.IsRepeated() != oField.IsRepeated() {
			return false
		}

		sType := sField.GetType()
		oType := oField.GetType()
		// Two fields with the same number do not have the same type
		if sType != oType {
			return false
		}
		if sType == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
			sFieldMessage := newMessageInternal(sField.GetMessageType())
			oFieldMessage := newMessageInternal(oField.GetMessageType())

			if !sFieldMessage.Compatible(oFieldMessage) {
				return false
			}
		}
	}
	return true
}

func (s *message) GetMessageField(name string) (Message, bool) {
	field := s.desc.FindFieldByName(name)
	if field == nil ||
		field.GetType() != descriptor.FieldDescriptorProto_TYPE_MESSAGE {

		return nil, false
	}
	return newMessageInternal(field.GetMessageType()), true
}
