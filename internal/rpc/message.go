package rpc

import (
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/util"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
)

// DynMessage wraps a grpc message while also allowing fields to be set.
type DynMessage interface {
	GrpcMsg() interface{}
	SetField(name string, val interface{})
	GetField(name string) (DynMessage, error)
}

type message struct {
	inner *dynamic.Message
}

func DynMessageFromProto(msg proto.Message) (DynMessage, error) {
	dyn, err := dynamic.AsDynamicMessage(msg)
	if err != nil {
		return nil, err
	}
	return &message{inner: dyn}, nil
}

func (m *message) GrpcMsg() interface{} {
	return m.inner
}

func (m *message) SetField(name string, val interface{}) {
	m.inner.SetFieldByName(name, val)
}

func (m *message) GetField(name string) (DynMessage, error) {
	f, err := m.inner.TryGetFieldByName(name)
	if err != nil {
		return nil, err
	}
	msg, ok := f.(proto.Message)
	if !ok {
		return nil, errdefs.InternalWithMsg("Field is not a message")
	}
	dyn, err := dynamic.AsDynamicMessage(msg)
	if err != nil {
		return nil, errdefs.InternalWithMsg(
			"convert proto msg to dynamic: %s",
			err,
		)
	}
	return &message{inner: dyn}, nil
}

// MessageDesc describes a grpc message
type MessageDesc interface {
	FullyQualifiedName() string
	// Compatible verifies if two Messages are compatible, meaning fields with
	// equal numbers have the same type.
	Compatible(other MessageDesc) bool
	// GetMessageField searches for a message field with the given name. It
	// returns a MessageDesc with the given field and true as the second value. If
	// no field with the given name exists, or the field has the wrong type, nil
	// and false is returned.
	GetMessageField(name string) (MessageDesc, bool)
	// NewEmpty returns an empty *dynamic.Message with the same fields as this
	// MessageDesc. The returned message can be used to call rpc methods that will
	// fill the fields.
	NewEmpty() DynMessage
}

type messageDesc struct {
	desc *desc.MessageDescriptor
}

func NewMessage(desc *desc.MessageDescriptor) (MessageDesc, error) {
	if ok, err := util.ArgNotNil(desc, "desc"); !ok {
		return nil, err
	}
	s := newMessageInternal(desc)
	return s, nil
}

func newMessageInternal(desc *desc.MessageDescriptor) MessageDesc {
	s := &messageDesc{desc: desc}
	return s
}

func (m *messageDesc) FullyQualifiedName() string {
	return m.desc.GetFullyQualifiedName()
}

func (m *messageDesc) Compatible(other MessageDesc) bool {
	otherMsg, ok := other.(*messageDesc)
	if !ok {
		return false
	}
	return m.cmpFields(otherMsg)
}

func (m *messageDesc) cmpFields(o *messageDesc) bool {
	for _, sField := range m.desc.GetFields() {
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

func (m *messageDesc) GetMessageField(name string) (MessageDesc, bool) {
	field := m.desc.FindFieldByName(name)
	if field == nil ||
		field.GetType() != descriptor.FieldDescriptorProto_TYPE_MESSAGE {

		return nil, false
	}
	return newMessageInternal(field.GetMessageType()), true
}

func (m *messageDesc) NewEmpty() DynMessage {
	return &message{inner: dynamic.NewMessage(m.desc)}
}
