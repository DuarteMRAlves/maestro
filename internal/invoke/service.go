package invoke

import (
	"fmt"
	"github.com/jhump/protoreflect/desc"
)

type Service interface {
	Methods() []MethodDescriptor
}

type service struct {
	methods []MethodDescriptor
}

func (s service) Methods() []MethodDescriptor {
	return s.methods
}

func newService(desc *desc.ServiceDescriptor) (Service, error) {
	methodDescriptors := desc.GetMethods()
	methods := make([]MethodDescriptor, 0, len(methodDescriptors))

	for _, d := range methodDescriptors {
		m, err := newMethodDescriptor(d)
		// Should never happen
		if err != nil {
			return nil, fmt.Errorf("create %s descriptor: %w", d.GetName(), err)
		}
		methods = append(methods, m)
	}

	return service{methods: methods}, nil
}
