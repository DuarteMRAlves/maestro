package invoke

import (
	"github.com/DuarteMRAlves/maestro/internal"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
)

type MessageGenerator func() DynamicMessage

type MessageDescriptor interface {
	Message() proto.Message
	MessageFields() map[internal.MessageField]MessageDescriptor
	MessageGenerator() MessageGenerator
}

type messageDescriptor struct {
	msg    proto.Message
	fields map[internal.MessageField]MessageDescriptor
	gen    MessageGenerator
}

func (d messageDescriptor) Message() proto.Message {
	return d.msg
}

func (d messageDescriptor) MessageFields() map[internal.MessageField]MessageDescriptor {
	return d.fields
}

func (d messageDescriptor) MessageGenerator() MessageGenerator {
	return d.gen
}

func NewMessageDescriptor(msg proto.Message) (MessageDescriptor, error) {
	d, err := desc.LoadMessageDescriptorForMessage(msg)
	if err != nil {
		return nil, err
	}
	return newMessageDescriptor(d)
}

func newMessageDescriptor(
	desc *desc.MessageDescriptor,
) (
	MessageDescriptor,
	error,
) {
	fields := map[internal.MessageField]MessageDescriptor{}
	for _, f := range desc.GetFields() {
		if f.GetType() == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
			var err error
			name := internal.NewMessageField(f.GetName())
			fields[name], err = newMessageDescriptor(f.GetMessageType())
			if err != nil {
				return nil, err
			}
		}
	}
	msg := dynamic.NewMessage(desc)
	gen := newGenFn(desc)
	return messageDescriptor{msg: msg, fields: fields, gen: gen}, nil
}

func newGenFn(desc *desc.MessageDescriptor) func() DynamicMessage {
	return func() DynamicMessage { return newDynamicMessageFromDesc(desc) }
}

func CompatibleDescriptors(d1, d2 MessageDescriptor) bool {
	desc1, err := desc.LoadMessageDescriptorForMessage(d1.Message())
	if err != nil {
		return false
	}
	desc2, err := desc.LoadMessageDescriptorForMessage(d2.Message())
	if err != nil {
		return false
	}
	return cmpFields(desc1, desc2)
}

func cmpFields(d1, d2 *desc.MessageDescriptor) bool {
	for _, f1 := range d1.GetFields() {
		number := f1.GetNumber()
		f2 := d2.FindFieldByNumber(number)

		// Ignore unmatched fields
		if f2 == nil {
			continue
		}

		// Both fields must be repeatable or not repeatable
		if f1.IsRepeated() != f2.IsRepeated() {
			return false
		}

		type1 := f1.GetType()
		type2 := f2.GetType()
		// Two fields with the same number do not have the same type
		if type1 != type2 {
			return false
		}
		if type1 == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
			msgField1 := f1.GetMessageType()
			msgField2 := f2.GetMessageType()

			if !cmpFields(msgField1, msgField2) {
				return false
			}
		}
	}
	return true
}
