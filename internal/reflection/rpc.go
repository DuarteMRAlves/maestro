package reflection

import (
	"github.com/DuarteMRAlves/maestro/internal/validate"
	"github.com/jhump/protoreflect/desc"
)

// RPC describes a rpc method of a given service.
type RPC interface {
	fullyQualifiedName() string
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

func (r *rpc) fullyQualifiedName() string {
	return r.desc.GetFullyQualifiedName()
}
