package grpc

import (
	"context"
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal"
	"github.com/DuarteMRAlves/maestro/internal/execute"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/grpcreflect"
	"google.golang.org/grpc"
	gr "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
	"google.golang.org/grpc/status"
	"time"
)

const reflectionServiceName = "grpc.reflection.v1alpha.ServerReflection"

var ReflectionMethodLoader execute.MethodLoader = &reflectionMethodLoader{}

type reflectionMethodLoader struct{}

func (m *reflectionMethodLoader) Load(methodCtx internal.MethodContext) (
	internal.UnaryMethod,
	error,
) {
	conn, err := grpc.Dial(methodCtx.Address().Unwrap(), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	services, err := listServices(ctx, conn)
	if err != nil {
		return nil, err
	}
	service, err := findService(services, methodCtx.Service())
	if err != nil {
		return nil, err
	}
	serviceDesc, err := resolveService(ctx, conn, service)
	if err != nil {
		return nil, err
	}
	return findMethod(serviceDesc.GetMethods(), methodCtx.Method())
}

func listServices(ctx context.Context, conn grpc.ClientConnInterface) (
	[]internal.Service,
	error,
) {
	stub := gr.NewServerReflectionClient(conn)
	c := grpcreflect.NewClient(ctx, stub)
	all, err := c.ListServices()
	if err != nil {
		st, _ := status.FromError(err)
		return nil, fmt.Errorf("list services: %w", st.Err())
	}
	// Filter the reflection service
	services := make([]internal.Service, 0, len(all)-1)
	for _, s := range all {
		if s != reflectionServiceName {
			services = append(services, internal.NewService(s))
		}
	}
	return services, nil
}

func findService(
	available []internal.Service,
	search internal.Service,
) (internal.Service, error) {
	if search.IsEmpty() {
		if len(available) == 1 {
			return available[0], nil
		}
		return internal.Service{}, notOneService
	} else {
		for _, s := range available {
			if search == s {
				return search, nil
			}
		}
		err := &internal.NotFound{Type: "service", Ident: search.Unwrap()}
		return internal.Service{}, err
	}
}

func resolveService(
	ctx context.Context,
	conn grpc.ClientConnInterface,
	service internal.Service,
) (*desc.ServiceDescriptor, error) {
	stub := gr.NewServerReflectionClient(conn)
	c := grpcreflect.NewClient(ctx, stub)
	descriptor, err := c.ResolveService(service.Unwrap())
	if err != nil {
		switch {
		case isGrpcErr(err):
			st, _ := status.FromError(err)
			err := st.Err()
			err = fmt.Errorf("resolve service %s: %w", service.Unwrap(), err)
			return nil, err
		case isElementNotFoundErr(err):
			err := &internal.NotFound{Type: "service", Ident: service.Unwrap()}
			return nil, fmt.Errorf("resolve service: %w", err)
		case isProtocolError(err):
			err := fmt.Errorf("resolve service %s: %w", service.Unwrap(), err)
			return nil, err
		default:
			// Should never happen as all errors should be caught by one
			// of the above options
			err := fmt.Errorf("resolve service %s: %w", service.Unwrap(), err)
			return nil, err
		}
	}
	return descriptor, nil
}

func isGrpcErr(err error) bool {
	_, ok := status.FromError(err)
	return ok
}

func isElementNotFoundErr(err error) bool {
	return grpcreflect.IsElementNotFoundError(err)
}

func isProtocolError(err error) bool {
	_, ok := err.(*grpcreflect.ProtocolError)
	return ok
}

func findMethod(
	available []*desc.MethodDescriptor,
	search internal.Method,
) (unaryMethod, error) {
	if search.IsEmpty() {
		if len(available) == 1 {
			return newUnaryMethodFromDescriptor(available[0]), nil
		}
		return unaryMethod{}, notOneMethod
	} else {
		for _, m := range available {
			if search.Unwrap() == m.GetName() {
				return newUnaryMethodFromDescriptor(m), nil
			}
		}
		err := &internal.NotFound{Type: "method", Ident: search.Unwrap()}
		return unaryMethod{}, err
	}
}
