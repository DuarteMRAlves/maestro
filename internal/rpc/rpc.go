package rpc

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/util"
	"github.com/jhump/protoreflect/desc"
)

// RPC describes a rpc method of a given service.
type RPC interface {
	Name() string
	FullyQualifiedName() string
	// InvokePath returns the path that should be used to invoke this rpc in a
	// remote server.
	InvokePath() string
	Service() Service
	Input() MessageDesc
	Output() MessageDesc
	IsUnary() bool
}

type rpc struct {
	desc *desc.MethodDescriptor
}

func newRPC(desc *desc.MethodDescriptor) (RPC, error) {
	if ok, err := util.ArgNotNil(desc, "desc"); !ok {
		return nil, err
	}
	r := &rpc{desc: desc}
	return r, nil
}

func (r *rpc) Name() string {
	return r.desc.GetName()
}

func (r *rpc) FullyQualifiedName() string {
	return r.desc.GetFullyQualifiedName()
}

func (r *rpc) InvokePath() string {
	return fmt.Sprintf("/%s/%s", r.Service().FullyQualifiedName(), r.Name())
}

func (r *rpc) Service() Service {
	return newServiceInternal(r.desc.GetService())
}

func (r *rpc) Input() MessageDesc {
	return newMessageInternal(r.desc.GetInputType())
}

func (r *rpc) Output() MessageDesc {
	return newMessageInternal(r.desc.GetOutputType())
}

func (r *rpc) IsUnary() bool {
	return !r.desc.IsServerStreaming() && !r.desc.IsClientStreaming()
}
