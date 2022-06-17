package grpc

import (
	"context"
	"fmt"
	"reflect"
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

func (m *ReflectionMethodLoader) Load(mid compiled.MethodID) (compiled.MethodDesc, error) {
	m.logger.Debugf("Load method with reflection: %#v", mid)

	mctx, ok := mid.(MethodID)
	if !ok {
		err := fmt.Errorf("unexpected type for method id: expected MethodContext, got %s", reflect.TypeOf(mid))
		return nil, err
	}

	conn, err := grpc.Dial(string(mctx.Address()), grpc.WithInsecure())
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
	service, err := findService(services, mctx.Service())
	if err != nil {
		return nil, err
	}
	serviceDesc, err := m.resolveService(ctx, conn, service)
	if err != nil {
		return nil, err
	}
	method, err := findMethod(serviceDesc.GetMethods(), mctx.Method())
	if err != nil {
		return nil, err
	}
	return newUnaryMethodFromDescriptor(method, mctx.Address().String()), nil
}

func (m *ReflectionMethodLoader) listServices(
	ctx context.Context, conn grpc.ClientConnInterface,
) ([]Service, error) {
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
	services := make([]Service, 0, len(all)-1)
	for _, s := range all {
		if s != reflectionServiceName {
			services = append(services, Service(s))
		}
	}
	return services, nil
}

func findService(available []Service, search Service) (Service, error) {
	if search.IsUnspecified() {
		if len(available) == 1 {
			return available[0], nil
		}
		return "", notOneService
	} else {
		for _, s := range available {
			if search == s {
				return search, nil
			}
		}
		return "", &serviceNotFound{srv: string(search)}
	}
}

func (m *ReflectionMethodLoader) resolveService(
	ctx context.Context, conn grpc.ClientConnInterface, service Service,
) (*desc.ServiceDescriptor, error) {
	var (
		descriptor *desc.ServiceDescriptor
		err        error
	)
	stub := gr.NewServerReflectionClient(conn)
	c := grpcreflect.NewClient(ctx, stub)

	retry.WhileTrue(func() bool {
		descriptor, err = c.ResolveService(string(service))
		st, _ := status.FromError(err)
		return st.Code() == codes.Unavailable
	}, &m.expBackoff)

	if err != nil {
		switch {
		case isGrpcErr(err):
			st, _ := status.FromError(err)
			err := st.Err()
			err = fmt.Errorf("resolve service %s: %w", service, err)
			return nil, err
		case isElementNotFoundErr(err):
			err := &serviceNotFound{srv: string(service)}
			return nil, fmt.Errorf("resolve service: %w", err)
		case isProtocolError(err):
			err := fmt.Errorf("resolve service %s: %w", service, err)
			return nil, err
		default:
			// Should never happen as all errors should be caught by one
			// of the above options
			err := fmt.Errorf("resolve service %s: %w", service, err)
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
	available []*desc.MethodDescriptor, search Method,
) (*desc.MethodDescriptor, error) {
	if search.IsUnspecified() {
		if len(available) == 1 {
			return available[0], nil
		}
		return nil, notOneMethod
	} else {
		for _, m := range available {
			if string(search) == m.GetName() {
				return m, nil
			}
		}
		return nil, &methodNotFound{meth: string(search)}
	}
}

type MethodID struct {
	address Address
	service Service
	method  Method
}

func (m MethodID) Address() Address { return m.address }

func (m MethodID) Service() Service { return m.service }

func (m MethodID) Method() Method { return m.method }

func (m MethodID) String() string {
	return fmt.Sprintf("%s/%s/%s", m.address, m.service, m.method)
}

func NewMethodID(
	address Address,
	service Service,
	method Method,
) MethodID {
	return MethodID{
		address: address,
		service: service,
		method:  method,
	}
}

// Address specifies the location of the server executing the
// stage method.
type Address string

func (a Address) IsEmpty() bool { return a == "" }

func (a Address) String() string {
	if a.IsEmpty() {
		return "*"
	}
	return string(a)
}

// Service specifies the name of the grpc service to execute.
type Service string

// IsUnspecified reports whether this service is either "" or "*".
func (s Service) IsUnspecified() bool { return s == "" || s == "*" }

func (s Service) String() string {
	if s.IsUnspecified() {
		return "*"
	}
	return string(s)
}

// Method specified the name of the grpc method to execute.
type Method string

// IsUnspecified reports whether this method is either "" or "*".
func (m Method) IsUnspecified() bool { return m == "" || m == "*" }

func (m Method) String() string {
	if m.IsUnspecified() {
		return "*"
	}
	return string(m)
}
