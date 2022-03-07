package invoke

import (
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
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
			err = errdefs.PrependMsg(err, "create service %s", desc.GetName())
			return nil, err
		}
		methods = append(methods, m)
	}

	return service{methods: methods}, nil
}
