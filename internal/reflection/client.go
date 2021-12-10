package reflection

import (
	"context"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/grpcreflect"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	gr "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
	"google.golang.org/grpc/status"
)

const reflectionServiceName = "grpc.reflection.v1alpha.ServerReflection"

// Client exposes an API for grpc reflection operations
type Client interface {
	ListServices() ([]string, error)
	ResolveService(name string) (*desc.ServiceDescriptor, error)
}

type client struct {
	client *grpcreflect.Client
}

// NewClient returns a new Client with the given context and connection
func NewClient(ctx context.Context, conn grpc.ClientConnInterface) Client {
	stub := gr.NewServerReflectionClient(conn)
	c := grpcreflect.NewClient(ctx, stub)
	return &client{client: c}
}

// ListServices lists the services available in the server. It does not show the
// reflection service that is activated.
func (c *client) ListServices() ([]string, error) {
	all, err := c.client.ListServices()
	if err != nil {
		return nil, handleGrpcError(err, "list services: ")
	}
	// Filter the reflection service
	services := make([]string, 0, len(all)-1)
	for _, r := range all {
		if r != reflectionServiceName {
			services = append(services, r)
		}
	}
	return services, nil
}

// ResolveService returns a descriptor for the service with the given name.
func (c *client) ResolveService(name string) (*desc.ServiceDescriptor, error) {
	descriptor, err := c.client.ResolveService(name)
	if err != nil {
		switch {
		case isGrpcErr(err):
			return nil, handleGrpcError(err, "resolve service: ")
		case isElementNotFoundErr(err):
			return nil, errdefs.NotFoundWithMsg(
				"resolve service: %v",
				err.Error())
		case isProtocolError(err):
			return nil, errdefs.UnknownWithError(err)
		default:
			// Should never happen as all errors should be caught by one
			// of the above options
			return nil, errdefs.InternalWithMsg("resolve service: %v", err)
		}
	}
	return descriptor, nil
}

func handleGrpcError(err error, prependMsg string) error {
	if err == nil {
		return nil
	}
	st, _ := status.FromError(err)
	switch st.Code() {
	case codes.Unavailable, codes.Unimplemented:
		// Unavailable is for the case where maestro is not running. When a
		// stage is not running, it is a failed precondition.
		// Unimplemented is when the maestro server does not implement a given
		// method. When a stage does not have an implemented method, in this
		// case reflection, it is a failed precondition.
		return errdefs.FailedPreconditionWithMsg("%v%v", prependMsg, st.Err())
	default:
		return errdefs.UnknownWithError(st.Err())
	}
}
