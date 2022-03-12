package invoke

import (
	"github.com/DuarteMRAlves/maestro/internal"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
)

// DynamicMessage offers a wrapper around a grpc message with some extra
// operations.
type DynamicMessage interface {
	GrpcMessage() proto.Message
	SetField(internal.MessageField, interface{}) error
	GetField(internal.MessageField) (DynamicMessage, error)
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
	field internal.MessageField,
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

func (dm *dynamicMessage) GetField(name internal.MessageField) (
	DynamicMessage,
	error,
) {
	field, err := dm.grpcMsg.TryGetFieldByName(name.Unwrap())
	if err != nil {
		return nil, err
	}
	msg, ok := field.(proto.Message)
	if !ok {
		return nil, errdefs.InternalWithMsg("Field is not a message")
	}
	dyn, err := NewDynamicMessage(msg)
	if err != nil {
		return nil, errdefs.InternalWithMsg("convert proto to dynamic: %s", err)
	}
	return dyn, nil
}
