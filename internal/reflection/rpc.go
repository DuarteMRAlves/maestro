package reflection

import (
	"github.com/DuarteMRAlves/maestro/internal/validate"
	"github.com/jhump/protoreflect/desc"
)

// RPC describes a rpc method of a given service.
type RPC interface {
	Name() string
	FullyQualifiedName() string
	Input() Message
	Output() Message
}

type rpc struct {
	desc *desc.MethodDescriptor
}

func newRPC(desc *desc.MethodDescriptor) (RPC, error) {
	if ok, err := validate.ArgNotNil(desc, "desc"); !ok {
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

func (r *rpc) Input() Message {
	return newMessageInternal(r.desc.GetInputType())
}

func (r *rpc) Output() Message {
	return newMessageInternal(r.desc.GetOutputType())
}
