package grpcw

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/DuarteMRAlves/maestro/internal/method"
	"github.com/DuarteMRAlves/maestro/internal/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/known/emptypb"
)

const reflectionServiceName = "grpc.reflection.v1alpha.ServerReflection"

var (
	errMalFormedAddress = errors.New("malformed address")
	errNotOneService    = errors.New("expected 1 available service")
	errNotOneMethod     = errors.New("expected 1 available method")
)

type Logger interface {
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
}

type ReflectionResolver struct {
	timeout    time.Duration
	expBackoff retry.ExponentialBackoff

	registry ProtoRegistry

	logger Logger
}

func NewReflectionResolver(
	timeout time.Duration, backoff retry.ExponentialBackoff, logger Logger,
) (*ReflectionResolver, error) {
	var registry ProtoRegistry
	registry = registry.RegisterFile(emptypb.File_google_protobuf_empty_proto)
	r := &ReflectionResolver{
		timeout:    timeout,
		expBackoff: backoff,
		registry:   registry,
		logger:     logger,
	}
	return r, nil
}

func (m *ReflectionResolver) Resolve(ctx context.Context, address string) (method.Desc, error) {
	m.logger.Debugf("Load method with reflection: %q\n", address)

	addr, err := m.parseAddress(address)
	if err != nil {
		return nil, err
	}

	conn, err := grpc.Dial(string(addr.Address()), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	ctx, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()

	services, err := m.listServices(ctx, conn)
	if err != nil {
		return nil, err
	}
	service, err := findService(services, addr.Service())
	if err != nil {
		return nil, err
	}
	serviceDesc, err := m.resolveService(ctx, conn, service)
	if err != nil {
		return nil, err
	}
	method, err := findMethod(serviceDesc.Methods(), addr.Method())
	if err != nil {
		return nil, err
	}
	return newUnaryMethodFromDescriptor(method, addr.Address().String()), nil
}

func (r *ReflectionResolver) parseAddress(address string) (addr, error) {
	var addr addr
	splits := strings.Split(address, "/")
	switch n := len(splits); n {
	// No backslash, only address was specified.
	case 1:
		addr.server = Address(splits[0])
	// Single backslash, this should divide the server address from the service.
	case 2:
		addr.server = Address(splits[0])
		addr.service = Service(splits[1])
	// Double backslash, this divides into server address, service and method.
	case 3:
		addr.server = Address(splits[0])
		addr.service = Service(splits[1])
		addr.method = Method(splits[2])
	default:
		return addr, errMalFormedAddress
	}
	return addr, nil
}

func (m *ReflectionResolver) listServices(
	ctx context.Context, conn grpc.ClientConnInterface,
) ([]Service, error) {
	var (
		all []string
		st  *status.Status
	)

	retry.WhileTrue(func() bool {
		stream, err := newBlockingReflectionStream(ctx, conn)
		if err != nil {
			st, _ = status.FromError(err)
			return st.Code() == codes.Unavailable
		}
		all, err = stream.listServiceNames()
		st, _ = status.FromError(err)
		return st.Code() == codes.Unavailable
	}, &m.expBackoff)
	if st.Err() != nil {
		return nil, fmt.Errorf("list services: %w", st.Err())
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
		return "", errNotOneService
	} else {
		for _, s := range available {
			if search == s {
				return search, nil
			}
		}
		return "", &serviceNotFound{srv: string(search)}
	}
}

func (m *ReflectionResolver) resolveService(
	ctx context.Context, conn grpc.ClientConnInterface, service Service,
) (protoreflect.ServiceDescriptor, error) {
	var (
		data [][]byte
		st   *status.Status
	)

	retry.WhileTrue(func() bool {
		stream, err := newBlockingReflectionStream(ctx, conn)
		if err != nil {
			st, _ = status.FromError(err)
			return st.Code() == codes.Unavailable
		}
		data, err = stream.filesForSymbol(string(service))
		st, _ = status.FromError(err)
		return st.Code() == codes.Unavailable
	}, &m.expBackoff)
	if st.Err() != nil {
		switch st.Code() {
		case codes.NotFound:
			return nil, &serviceNotFound{srv: string(service)}
		default:
			return nil, fmt.Errorf("resolve service %s: %w", service, st.Err())
		}
	}
	if err := m.registerFiles(data); err != nil {
		return nil, fmt.Errorf("resolve service %s: %w", service, err)
	}
	d, err := m.registry.FindDescriptorByName(protoreflect.FullName(service))
	if err != nil {
		return nil, err
	}
	desc, ok := d.(protoreflect.ServiceDescriptor)
	if !ok {
		return nil, &notService{symb: string(service)}
	}
	return desc, nil
}

func (m *ReflectionResolver) registerFiles(data [][]byte) error {
	for _, buf := range data {
		var descPb descriptorpb.FileDescriptorProto
		if err := proto.Unmarshal(buf, &descPb); err != nil {
			return err
		}
		desc, err := protodesc.NewFile(&descPb, m.registry)
		if err != nil {
			return err
		}
		m.registry = m.registry.RegisterFile(desc)
	}
	return nil
}

func findMethod(
	available protoreflect.MethodDescriptors, search Method,
) (protoreflect.MethodDescriptor, error) {
	if search.IsUnspecified() {
		if available.Len() == 1 {
			return available.Get(0), nil
		}
		return nil, errNotOneMethod
	} else {
		m := available.ByName(protoreflect.Name(search))
		if m == nil {
			return nil, &methodNotFound{meth: string(search)}
		}
		return m, nil
	}
}

type addr struct {
	server  Address
	service Service
	method  Method
}

func (m addr) Address() Address { return m.server }

func (m addr) Service() Service { return m.service }

func (m addr) Method() Method { return m.method }

func (m addr) String() string {
	return fmt.Sprintf("%s/%s/%s", m.server, m.service, m.method)
}

func NewAddress(
	address Address,
	service Service,
	method Method,
) addr {
	return addr{
		server:  address,
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

type serviceNotFound struct {
	srv string
}

func (err *serviceNotFound) NotFound() {}

func (err *serviceNotFound) Error() string {
	return fmt.Sprintf("service not found: %s", err.srv)
}

type notService struct {
	symb string
}

func (err *notService) Error() string {
	return fmt.Sprintf("symbol not a service: %q", err.symb)
}

type methodNotFound struct {
	meth string
}

func (err *methodNotFound) NotFound() {}

func (err *methodNotFound) Error() string {
	return fmt.Sprintf("method not found: %s", err.meth)
}
