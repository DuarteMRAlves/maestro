package reflection

import (
	"github.com/DuarteMRAlves/maestro/internal/validate"
	"github.com/jhump/protoreflect/desc"
)

// Service describes a grpc service
type Service interface {
	fullyQualifiedName() string
	RPCs() []RPC
}

type service struct {
	desc *desc.ServiceDescriptor
}

func newService(desc *desc.ServiceDescriptor) (Service, error) {
	if ok, err := validate.ArgNotNil(desc, "desc"); !ok {
		return nil, err
	}
	s := &service{desc: desc}
	return s, nil
}

func (s *service) fullyQualifiedName() string {
	return s.desc.GetFullyQualifiedName()
}

func (s *service) RPCs() []RPC {
	methodDescriptors := s.desc.GetMethods()
	rpcs := make([]RPC, 0, len(methodDescriptors))

	for _, m := range methodDescriptors {
		rpc, err := newRPC(m)
		// Should never happen
		if err != nil {
			panic(err)
		}
		rpcs = append(rpcs, rpc)
	}

	return rpcs
}
