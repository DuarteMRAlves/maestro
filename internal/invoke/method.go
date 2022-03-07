package invoke

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/jhump/protoreflect/desc"
)

type FullMethod interface {
	FullMethod()
	Unwrap() string
}

type MethodDescriptor interface {
	FullMethod() FullMethod
	Input() MessageDescriptor
	Output() MessageDescriptor
	IsUnary() bool
}

type fullMethod string

func (m fullMethod) FullMethod() {}

func (m fullMethod) Unwrap() string { return string(m) }

func newFullMethod(desc *desc.MethodDescriptor) FullMethod {
	m := fmt.Sprintf(
		"/%s/%s",
		desc.GetService().GetFullyQualifiedName(),
		desc.GetName(),
	)
	return fullMethod(m)
}

type methodDescriptor struct {
	method  FullMethod
	input   MessageDescriptor
	output  MessageDescriptor
	isUnary bool
}

func (d methodDescriptor) FullMethod() FullMethod {
	return d.method
}

func (d methodDescriptor) Input() MessageDescriptor {
	return d.input
}

func (d methodDescriptor) Output() MessageDescriptor {
	return d.output
}

func (d methodDescriptor) IsUnary() bool {
	return d.isUnary
}

func newMethodDescriptor(desc *desc.MethodDescriptor) (
	MethodDescriptor,
	error,
) {
	method := newFullMethod(desc)
	input, err := newMessageDescriptor(desc.GetInputType())
	if err != nil {
		return nil, errdefs.PrependMsg(err, "create method")
	}
	output, err := newMessageDescriptor(desc.GetOutputType())
	if err != nil {
		return nil, errdefs.PrependMsg(err, "create method")
	}
	isUnary := !desc.IsServerStreaming() && !desc.IsClientStreaming()
	methodDesc := methodDescriptor{
		method:  method,
		input:   input,
		output:  output,
		isUnary: isUnary,
	}
	return methodDesc, nil
}
