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
	clientBuilder compiled.UnaryClientBuilder
	input         messageDescriptor
	output        messageDescriptor
}

func (d unaryMethod) ClientBuilder() compiled.UnaryClientBuilder {
	return d.clientBuilder
}

func (d unaryMethod) Input() compiled.MessageDesc {
	return d.input
}

func (d unaryMethod) Output() compiled.MessageDesc {
	return d.output
}

func newUnaryMethod(
	invokePath string,
	inDesc, outDesc messageDescriptor,
) unaryMethod {
	return unaryMethod{
		input:         inDesc,
		output:        outDesc,
		clientBuilder: newClientBuilder(invokePath, outDesc.EmptyGen()),
	}
}

func newUnaryMethodFromDescriptor(desc *desc.MethodDescriptor) unaryMethod {
	invokePath := methodInvokePath(desc)
	input := messageDescriptor{desc: desc.GetInputType()}
	output := messageDescriptor{desc: desc.GetOutputType()}

	return newUnaryMethod(invokePath, input, output)
}

func newClientBuilder(
	invokePath string,
	emptyGen compiled.EmptyMessageGen,
) compiled.UnaryClientBuilder {
	return func(address compiled.Address) (compiled.Conn, error) {
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
