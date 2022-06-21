package grpc

import (
	"fmt"

	"github.com/DuarteMRAlves/maestro/internal/message"
	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
)

type messageInstance struct {
	dynMsg *dynamic.Message
}

func newMessage(msg proto.Message) (*messageInstance, error) {
	grpcMsg, err := dynamic.AsDynamicMessage(msg)
	if err != nil {
		return nil, fmt.Errorf("convert proto to dynamic: %w", err)
	}
	return &messageInstance{dynMsg: grpcMsg}, nil
}

func newMessageFromDescriptor(desc *desc.MessageDescriptor) *messageInstance {
	return &messageInstance{dynMsg: dynamic.NewMessage(desc)}
}

func (dm *messageInstance) Set(field message.Field, msg message.Instance) error {
	grpcMsg, ok := msg.(*messageInstance)
	if !ok {
		return notGrpcMessage
	}
	return dm.dynMsg.TrySetFieldByName(string(field), grpcMsg.dynMsg)
}

func (dm *messageInstance) Get(name message.Field) (message.Instance, error) {
	field, err := dm.dynMsg.TryGetFieldByName(string(name))
	if err != nil {
		return nil, err
	}
	msg, ok := field.(proto.Message)
	if !ok {
		return nil, &fieldNotMessage{
			MsgType: dm.dynMsg.XXX_MessageName(),
			Field:   string(name),
		}
	}
	dyn, err := newMessage(msg)
	if err != nil {
		return nil, err
	}
	return dyn, nil
}
