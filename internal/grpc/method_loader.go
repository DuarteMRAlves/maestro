package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/DuarteMRAlves/maestro/internal/compiled"
	"github.com/DuarteMRAlves/maestro/internal/retry"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/grpcreflect"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	gr "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
	"google.golang.org/grpc/status"
)

const reflectionServiceName = "grpc.reflection.v1alpha.ServerReflection"

type Logger interface {
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
}

type ReflectionMethodLoader struct {
	timeout    time.Duration
	expBackoff retry.ExponentialBackoff

	logger Logger
}

func NewReflectionMethodLoader(
	timeout time.Duration, backoff retry.ExponentialBackoff, logger Logger,
) *ReflectionMethodLoader {
	return &ReflectionMethodLoader{
		timeout:    timeout,
		expBackoff: backoff,
		logger:     logger,
	}
}

func (m *ReflectionMethodLoader) Load(methodCtx *compiled.MethodContext) (
	compiled.UnaryMethod,
	error,
) {
	m.logger.Debugf("Load method with reflection: %#v", methodCtx)
	conn, err := grpc.Dial(methodCtx.Address().Unwrap(), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	services, err := m.listServices(ctx, conn)
	if err != nil {
		return nil, err
	}
	service, err := findService(services, methodCtx.Service())
	if err != nil {
		return nil, err
	}
	serviceDesc, err := m.resolveService(ctx, conn, service)
	if err != nil {
		return nil, err
	}
	return findMethod(serviceDesc.GetMethods(), methodCtx.Method())
}

func (m *ReflectionMethodLoader) listServices(
	ctx context.Context, conn grpc.ClientConnInterface,
) ([]compiled.Service, error) {
	var (
		all []string
		err error
	)
	stub := gr.NewServerReflectionClient(conn)
	c := grpcreflect.NewClient(ctx, stub)

	retry.WhileTrue(func() bool {
		all, err = c.ListServices()
		st, _ := status.FromError(err)
		err = st.Err()
		return st.Code() == codes.Unavailable
	}, &m.expBackoff)
	if err != nil {
		return nil, fmt.Errorf("list services: %w", err)
	}
	// Filter the reflection service
	services := make([]compiled.Service, 0, len(all)-1)
	for _, s := range all {
		if s != reflectionServiceName {
			services = append(services, compiled.NewService(s))
		}
	}
	return services, nil
}

func findService(
	available []compiled.Service,
	search compiled.Service,
) (compiled.Service, error) {
	if search.IsEmpty() {
		if len(available) == 1 {
			return available[0], nil
		}
		return compiled.Service{}, notOneService
	} else {
		for _, s := range available {
			if search == s {
				return search, nil
			}
		}
		return compiled.Service{}, &serviceNotFound{srv: search.Unwrap()}
	}
}

func (m *ReflectionMethodLoader) resolveService(
	ctx context.Context,
	conn grpc.ClientConnInterface,
	service compiled.Service,
) (*desc.ServiceDescriptor, error) {
	var (
		descriptor *desc.ServiceDescriptor
		err        error
	)
	stub := gr.NewServerReflectionClient(conn)
	c := grpcreflect.NewClient(ctx, stub)

	retry.WhileTrue(func() bool {
		descriptor, err = c.ResolveService(service.Unwrap())
		st, _ := status.FromError(err)
		return st.Code() == codes.Unavailable
	}, &m.expBackoff)

	if err != nil {
		switch {
		case isGrpcErr(err):
			st, _ := status.FromError(err)
			err := st.Err()
			err = fmt.Errorf("resolve service %s: %w", service.Unwrap(), err)
			return nil, err
		case isElementNotFoundErr(err):
			err := &serviceNotFound{srv: service.Unwrap()}
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
	search compiled.Method,
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
		return unaryMethod{}, &methodNotFound{meth: search.Unwrap()}
	}
}
