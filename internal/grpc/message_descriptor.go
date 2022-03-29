package grpc

import (
	"github.com/DuarteMRAlves/maestro/internal"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jhump/protoreflect/desc"
)

type messageDescriptor struct {
	desc *desc.MessageDescriptor
}

func newMessageDescriptor(msg proto.Message) (messageDescriptor, error) {
	d, err := desc.LoadMessageDescriptorForMessage(msg)
	if err != nil {
		return messageDescriptor{}, err
	}
	return messageDescriptor{desc: d}, nil
}

func (d messageDescriptor) EmptyGen() internal.EmptyMessageGen {
	return func() internal.Message { return newMessageFromDescriptor(d.desc) }
}

func (d messageDescriptor) GetField(name internal.MessageField) (
	internal.MessageDesc,
	error,
) {
	field := d.desc.FindFieldByName(name.Unwrap())
	if field == nil {
		err := &fieldNotFound{
			msgType: d.desc.GetFullyQualifiedName(), field: name.Unwrap(),
		}
		return nil, err
	}
	if field.GetType() != descriptor.FieldDescriptorProto_TYPE_MESSAGE {
		return nil, &fieldNotMessage{
			MsgType: d.desc.GetFullyQualifiedName(),
			Field:   name.Unwrap(),
		}
	}
	return messageDescriptor{desc: field.GetMessageType()}, nil
}

func (d messageDescriptor) Compatible(other internal.MessageDesc) bool {
	grpcOther, ok := other.(messageDescriptor)
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
