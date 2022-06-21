package grpc

import (
	"context"
	"fmt"

	"github.com/DuarteMRAlves/maestro/internal/message"
	"github.com/DuarteMRAlves/maestro/internal/method"
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
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

func newUnaryMethodFromDescriptor(desc *desc.MethodDescriptor, address string) unaryMethod {
	invokePath := methodInvokePath(desc)
	input := messageType{desc: desc.GetInputType()}
	output := messageType{desc: desc.GetOutputType()}

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

func methodInvokePath(desc *desc.MethodDescriptor) string {
	return fmt.Sprintf(
		"/%s/%s",
		desc.GetService().GetFullyQualifiedName(),
		desc.GetName(),
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

	grpcReq, ok := req.(*messageInstance)
	if !ok {
		return nil, notGrpcMessage
	}
	grpcRep, ok := rep.(*messageInstance)
	if !ok {
		return nil, notGrpcMessage
	}

	err := c.conn.Invoke(ctx, c.invokePath, grpcReq.dynMsg, grpcRep.dynMsg)
	if err != nil {
		st, _ := status.FromError(err)
		return nil, fmt.Errorf("invoke %s: %w", c.invokePath, st.Err())
	}
	return grpcRep, nil
}

func (c unaryClient) Close() error {
	return c.conn.Close()
}
