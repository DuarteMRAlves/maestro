package grpcw

import (
	"fmt"

	"github.com/DuarteMRAlves/maestro/internal/message"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

type messageInstance struct {
	m protoreflect.Message
}

func (mi messageInstance) Set(field message.Field, value message.Instance) error {
	i, ok := value.(messageInstance)
	if !ok {
		return fmt.Errorf("value not of type messageInstance: %v", value)
	}
	fd := mi.searchFieldDescriptor(field)
	if fd == nil {
		return &errUnknownField{Field: field}
	}
	v := protoreflect.ValueOfMessage(i.m)
	mi.m.Set(fd, v)
	return nil
}

func (mi messageInstance) Get(field message.Field) (message.Instance, error) {
	fd := mi.searchFieldDescriptor(field)
	if fd == nil {
		return nil, &errUnknownField{Field: field}
	}
	v := mi.m.Get(fd).Interface()
	switch x := v.(type) {
	case protoreflect.Message:
		return messageInstance{m: x}, nil
	default:
		return nil, fmt.Errorf("field type not message: %q", field)
	}
}

func (mi messageInstance) searchFieldDescriptor(field message.Field) protoreflect.FieldDescriptor {
	fields := mi.m.Descriptor().Fields()
	return fields.ByName(protoreflect.Name(field))
}

type messageType struct {
	t protoreflect.MessageType
}

func (t messageType) Build() message.Instance {
	return messageInstance{t.t.New()}
}

func (t messageType) Subfield(field message.Field) (message.Type, error) {
	fd := t.searchFieldDescriptor(field)
	if fd == nil {
		return nil, &errUnknownField{Field: field}
	}
	if fd.Kind() != protoreflect.MessageKind {
		return nil, &fieldNotMessageKind{
			MsgType: string(t.t.Descriptor().Name()),
			Field:   string(field),
		}
	}
	typ := dynamicpb.NewMessageType(fd.Message())
	return messageType{typ}, nil
}

func (t messageType) Compatible(o message.Type) bool {
	other, ok := o.(messageType)
	if !ok {
		return false
	}
	return compatibleMessageDescriptors(t.t.Descriptor(), other.t.Descriptor())
}

func compatibleMessageDescriptors(d1, d2 protoreflect.MessageDescriptor) bool {
	fields1 := d1.Fields()
	fields2 := d2.Fields()
	for i1 := 0; i1 < fields1.Len(); i1++ {
		f1 := fields1.Get(i1)
		f2 := fields2.ByNumber(f1.Number())

		// Ignore unmatched fields
		if f2 == nil {
			continue
		}

		// Both fields must have the same cardinality
		if f1.Cardinality() != f2.Cardinality() {
			return false
		}

		// Fields with the same number must have the same kind
		if f1.Kind() != f2.Kind() {
			return false
		}

		// If the fields are messages, they must also be compatible
		if f1.Kind() == protoreflect.MessageKind {
			if !compatibleMessageDescriptors(f1.Message(), f2.Message()) {
				return false
			}
		}
	}
	return true
}

func (t messageType) searchFieldDescriptor(field message.Field) protoreflect.FieldDescriptor {
	fields := t.t.Descriptor().Fields()
	return fields.ByName(protoreflect.Name(field))
}

type errUnknownField struct {
	Field message.Field
}

func (err *errUnknownField) Error() string {
	return fmt.Sprintf("unknown field: %q", err.Field)
}

type fieldNotMessageKind struct {
	MsgType string
	Field   string
}

func (err *fieldNotMessageKind) Error() string {
	format := "kind for field %q of %q is not a message"
	return fmt.Sprintf(format, err.Field, err.MsgType)
}
