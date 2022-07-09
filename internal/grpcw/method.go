package grpcw

import (
	"context"
	"errors"
	"fmt"

	"github.com/DuarteMRAlves/maestro/internal/message"
	"github.com/DuarteMRAlves/maestro/internal/method"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

var (
	errNotGrpcMessage = errors.New("message is not grpc")
)

type unaryMethod struct {
	dialer method.DialFunc
	input  messageType
	output messageType
}

func (d unaryMethod) Dial() (method.Conn, error) {
	return d.dialer.Dial()
}

func (d unaryMethod) Input() message.Type {
	return d.input
}

func (d unaryMethod) Output() message.Type {
	return d.output
}

func newUnaryMethod(
	address string,
	invokePath string,
	inDesc, outDesc messageType,
) unaryMethod {
	return unaryMethod{
		input:  inDesc,
		output: outDesc,
		dialer: newDialFunc(address, invokePath, outDesc.Build),
	}
}

func newUnaryMethodFromDescriptor(desc protoreflect.MethodDescriptor, address string) unaryMethod {
	invokePath := methodInvokePath(desc)
	input := messageType{t: dynamicpb.NewMessageType(desc.Input())}
	output := messageType{t: dynamicpb.NewMessageType(desc.Output())}

	return newUnaryMethod(address, invokePath, input, output)
}

func newDialFunc(
	address string,
	invokePath string,
	emptyGen message.BuildFunc,
) method.DialFunc {
	return func() (method.Conn, error) {
		conn, err := grpc.Dial(string(address), grpc.WithInsecure())
		if err != nil {
			return nil, err
		}
		return unaryClient{
			conn:       conn,
			invokePath: invokePath,
			buildFunc:  emptyGen,
		}, nil
	}
}

func methodInvokePath(desc protoreflect.MethodDescriptor) string {
	return fmt.Sprintf(
		"/%s/%s",
		desc.FullName().Parent(),
		desc.FullName().Name(),
	)
}

type unaryClient struct {
	conn       *grpc.ClientConn
	invokePath string
	buildFunc  message.BuildFunc
}

func (c unaryClient) Call(
	ctx context.Context,
	req message.Instance,
) (message.Instance, error) {
	rep := c.buildFunc()

	reqInst, ok := req.(messageInstance)
	if !ok {
		return nil, errNotGrpcMessage
	}
	repInst, ok := rep.(messageInstance)
	if !ok {
		return nil, errNotGrpcMessage
	}

	err := c.conn.Invoke(ctx, c.invokePath, reqInst.m.Interface(), repInst.m.Interface())
	if err != nil {
		st, _ := status.FromError(err)
		return nil, fmt.Errorf("invoke %s: %w", c.invokePath, st.Err())
	}
	return repInst, nil
}

func (c unaryClient) Close() error {
	return c.conn.Close()
}
