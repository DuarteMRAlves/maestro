package grpc

import (
	"github.com/DuarteMRAlves/maestro/internal/message"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jhump/protoreflect/desc"
)

type messageType struct {
	desc *desc.MessageDescriptor
}

func newMessageDescriptor(msg proto.Message) (messageType, error) {
	d, err := desc.LoadMessageDescriptorForMessage(msg)
	if err != nil {
		return messageType{}, err
	}
	return messageType{desc: d}, nil
}

func (d messageType) Build() message.Instance {
	return newMessageFromDescriptor(d.desc)
}

func (d messageType) Subfield(name message.Field) (
	message.Type,
	error,
) {
	field := d.desc.FindFieldByName(string(name))
	if field == nil {
		err := &fieldNotFound{
			msgType: d.desc.GetFullyQualifiedName(), field: string(name),
		}
		return nil, err
	}
	if field.GetType() != descriptor.FieldDescriptorProto_TYPE_MESSAGE {
		return nil, &fieldNotMessage{
			MsgType: d.desc.GetFullyQualifiedName(),
			Field:   string(name),
		}
	}
	return messageType{desc: field.GetMessageType()}, nil
}

func (d messageType) Compatible(other message.Type) bool {
	grpcOther, ok := other.(messageType)
	if !ok {
		return false
	}
	return equalDescriptors(d.desc, grpcOther.desc)
}

func equalDescriptors(d1, d2 *desc.MessageDescriptor) bool {
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

			if !equalDescriptors(msgField1, msgField2) {
				return false
			}
		}
	}
	return true
}
