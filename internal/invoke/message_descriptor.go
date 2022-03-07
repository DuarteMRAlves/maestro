package invoke

import (
	"github.com/DuarteMRAlves/maestro/internal/domain"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
)

type MessageGenerator func() DynamicMessage

type MessageDescriptor interface {
	MessageFields() map[domain.MessageField]MessageDescriptor
	MessageGenerator() MessageGenerator
}

type messageDescriptor struct {
	desc   *desc.MessageDescriptor
	fields map[domain.MessageField]MessageDescriptor
}

func (d messageDescriptor) MessageFields() map[domain.MessageField]MessageDescriptor {
	return d.fields
}

func (d messageDescriptor) MessageGenerator() MessageGenerator {
	return func() DynamicMessage {
		return newDynamicMessage(dynamic.NewMessage(d.desc))
	}
}

func newMessageDescriptor(desc *desc.MessageDescriptor) (
	MessageDescriptor,
	error,
) {
	fields := map[domain.MessageField]MessageDescriptor{}
	for _, f := range desc.GetFields() {
		if f.GetType() == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
			name, err := domain.NewMessageField(f.GetName())
			// Should never happen as the name should be valid
			if err != nil {
				return nil, err
			}
			fields[name], err = newMessageDescriptor(f.GetMessageType())
			if err != nil {
				return nil, err
			}
		}
	}
	return messageDescriptor{desc: desc, fields: fields}, nil
}

func CompatibleDescriptors(d1, d2 MessageDescriptor) bool {
	msgDesc1, ok := d1.(messageDescriptor)
	if !ok {
		return false
	}
	msgDesc2, ok := d2.(messageDescriptor)
	if !ok {
		return false
	}
	return cmpFields(msgDesc1.desc, msgDesc2.desc)
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
