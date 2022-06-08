package grpc

import (
	"fmt"

	"github.com/DuarteMRAlves/maestro/internal/compiled"
	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
)

type message struct {
	dynMsg *dynamic.Message
}

func newMessage(msg proto.Message) (*message, error) {
	grpcMsg, err := dynamic.AsDynamicMessage(msg)
	if err != nil {
		return nil, fmt.Errorf("convert proto to dynamic: %w", err)
	}
	return &message{dynMsg: grpcMsg}, nil
}

func newMessageFromDescriptor(desc *desc.MessageDescriptor) *message {
	return &message{dynMsg: dynamic.NewMessage(desc)}
}

func (dm *message) SetField(
	field compiled.MessageField,
	msg compiled.Message,
) error {
	grpcMsg, ok := msg.(*message)
	if !ok {
		return notGrpcMessage
	}
	return dm.dynMsg.TrySetFieldByName(field.Unwrap(), grpcMsg.dynMsg)
}

func (dm *message) GetField(name compiled.MessageField) (
	compiled.Message,
	error,
) {
	field, err := dm.dynMsg.TryGetFieldByName(name.Unwrap())
	if err != nil {
		return nil, err
	}
	msg, ok := field.(proto.Message)
	if !ok {
		return nil, &fieldNotMessage{
			MsgType: dm.dynMsg.XXX_MessageName(),
			Field:   name.Unwrap(),
		}
	}
	dyn, err := newMessage(msg)
	if err != nil {
		return nil, err
	}
	return dyn, nil
}
