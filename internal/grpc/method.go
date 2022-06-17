package grpc

import (
	"context"
	"fmt"

	"github.com/DuarteMRAlves/maestro/internal/compiled"
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type unaryMethod struct {
	dialer compiled.DialFunc
	input  messageDescriptor
	output messageDescriptor
}

func (d unaryMethod) Dial() (compiled.Conn, error) {
	return d.dialer.Dial()
}

func (d unaryMethod) Input() compiled.MessageDesc {
	return d.input
}

func (d unaryMethod) Output() compiled.MessageDesc {
	return d.output
}

func newUnaryMethod(
	address string,
	invokePath string,
	inDesc, outDesc messageDescriptor,
) unaryMethod {
	return unaryMethod{
		input:  inDesc,
		output: outDesc,
		dialer: newDialFunc(address, invokePath, outDesc.EmptyGen()),
	}
}

func newUnaryMethodFromDescriptor(desc *desc.MethodDescriptor, address string) unaryMethod {
	invokePath := methodInvokePath(desc)
	input := messageDescriptor{desc: desc.GetInputType()}
	output := messageDescriptor{desc: desc.GetOutputType()}

	return newUnaryMethod(address, invokePath, input, output)
}

func newDialFunc(
	address string,
	invokePath string,
	emptyGen compiled.EmptyMessageGen,
) compiled.DialFunc {
	return func() (compiled.Conn, error) {
		conn, err := grpc.Dial(string(address), grpc.WithInsecure())
		if err != nil {
			return nil, err
		}
		return unaryClient{
			conn:       conn,
			invokePath: invokePath,
			emptyGen:   emptyGen,
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
	emptyGen   compiled.EmptyMessageGen
}

func (c unaryClient) Call(
	ctx context.Context,
	req compiled.Message,
) (compiled.Message, error) {
	rep := c.emptyGen()

	grpcReq, ok := req.(*message)
	if !ok {
		return nil, notGrpcMessage
	}
	grpcRep, ok := rep.(*message)
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
