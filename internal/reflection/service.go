package reflection

import (
	"github.com/DuarteMRAlves/maestro/internal/util"
	"github.com/jhump/protoreflect/desc"
)

// Service describes a grpc service
type Service interface {
	Name() string
	FullyQualifiedName() string
	RPCs() []RPC
}

type service struct {
	desc *desc.ServiceDescriptor
}

func newService(desc *desc.ServiceDescriptor) (Service, error) {
	if ok, err := util.ArgNotNil(desc, "desc"); !ok {
		return nil, err
	}
	s := newServiceInternal(desc)
	return s, nil
}

func newServiceInternal(desc *desc.ServiceDescriptor) Service {
	return &service{desc: desc}
}

func (s *service) Name() string {
	return s.desc.GetName()
}

func (s *service) FullyQualifiedName() string {
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
