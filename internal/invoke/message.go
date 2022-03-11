package invoke

import (
	"github.com/DuarteMRAlves/maestro/internal/domain"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
)

type FieldGetter func(DynamicMessage) DynamicMessageResult

// DynamicMessage offers a wrapper around a grpc message with some extra
// operations.
type DynamicMessage interface {
	GrpcMessage() proto.Message
	SetField(domain.MessageField, interface{}) error
}

type dynamicMessage struct {
	grpcMsg *dynamic.Message
}

func NewDynamicMessage(msg proto.Message) (DynamicMessage, error) {
	grpcMsg, err := dynamic.AsDynamicMessage(msg)
	if err != nil {
		return nil, err
	}
	return &dynamicMessage{grpcMsg: grpcMsg}, nil
}

func newDynamicMessageFromDesc(desc *desc.MessageDescriptor) DynamicMessage {
	return &dynamicMessage{grpcMsg: dynamic.NewMessage(desc)}
}

func (dm *dynamicMessage) GrpcMessage() proto.Message {
	return dm.grpcMsg
}

func (dm *dynamicMessage) SetField(
	field domain.MessageField,
	val interface{},
) error {
	msg := dynamic.NewMessage(dm.grpcMsg.GetMessageDescriptor())
	msg.Merge(dm.GrpcMessage())
	if err := msg.TrySetFieldByName(field.Unwrap(), val); err != nil {
		return err
	}
	dm.grpcMsg = msg
	return nil
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
		dyn, err := NewDynamicMessage(msg)
		if err != nil {
			err = errdefs.InternalWithMsg("convert proto to dynamic: %s", err)
			return ErrDynamicMessage(err)
		}
		return SomeDynamicMessage(dyn)
	}
}
