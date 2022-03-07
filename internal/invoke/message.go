package invoke

import (
	"github.com/DuarteMRAlves/maestro/internal/domain"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/dynamic"
)

type FieldSetter func(DynamicMessage, interface{}) DynamicMessageResult
type FieldGetter func(DynamicMessage) DynamicMessageResult

// DynamicMessage offers a wrapper around a grpc message with some extra
// operations.
type DynamicMessage interface {
	GrpcMessage() proto.Message
}

type dynamicMessage struct {
	grpcMsg *dynamic.Message
}

func newDynamicMessage(msg *dynamic.Message) DynamicMessage {
	return dynamicMessage{grpcMsg: msg}
}

func (dm dynamicMessage) GrpcMessage() proto.Message {
	return dm.grpcMsg
}

func NewFieldSetter(field domain.MessageField) FieldSetter {
	return func(m DynamicMessage, val interface{}) DynamicMessageResult {
		old, err := dynamic.AsDynamicMessage(m.GrpcMessage())
		if err != nil {
			return ErrDynamicMessage(err)
		}
		msg := dynamic.NewMessage(old.GetMessageDescriptor())
		msg.Merge(m.GrpcMessage())
		msg.SetFieldByName(field.Unwrap(), val)
		return SomeDynamicMessage(newDynamicMessage(msg))
	}
}

func NewFieldGetter(field domain.MessageField) FieldGetter {
	return func(m DynamicMessage) DynamicMessageResult {
		old, err := dynamic.AsDynamicMessage(m.GrpcMessage())
		if err != nil {
			return ErrDynamicMessage(err)
		}
		f, err := old.TryGetFieldByName(field.Unwrap())
		if err != nil {
			return ErrDynamicMessage(err)
		}
		msg, ok := f.(proto.Message)
		if !ok {
			err = errdefs.InternalWithMsg("Field is not a message")
			return ErrDynamicMessage(err)
		}
		dyn, err := dynamic.AsDynamicMessage(msg)
		if err != nil {
			err = errdefs.InternalWithMsg("convert proto to dynamic: %s", err)
			return ErrDynamicMessage(err)
		}
		return SomeDynamicMessage(newDynamicMessage(dyn))
	}
}
