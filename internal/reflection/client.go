package reflection

import (
	"context"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
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
		return nil, handleGrpcError(err)
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

func handleGrpcError(err error) error {
	if err == nil {
		return nil
	}
	st, _ := status.FromError(err)
	switch st.Code() {
	case codes.Unavailable:
		return errdefs.UnavailableWithMsg("list services: %v", st.Err())
	case codes.Unimplemented:
		return errdefs.FailedPreconditionWithMsg("list services: %v", st.Err())
	default:
		return errdefs.UnknownWithError(st.Err())
	}
}
