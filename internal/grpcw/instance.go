package grpcw

import (
	"fmt"

	"github.com/DuarteMRAlves/maestro/internal/message"
	"google.golang.org/protobuf/reflect/protoreflect"
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

type errUnknownField struct {
	Field message.Field
}

func (err *errUnknownField) Error() string {
	return fmt.Sprintf("unknown field: %q", err.Field)
}
