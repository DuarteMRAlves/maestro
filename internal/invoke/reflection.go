package invoke

import (
	"context"
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/jhump/protoreflect/grpcreflect"
	"google.golang.org/grpc"
	gr "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
	"google.golang.org/grpc/status"
)

const reflectionServiceName = "grpc.reflection.v1alpha.ServerReflection"

func listServices(conn grpc.ClientConnInterface) func(context.Context) (
	[]internal.Service,
	error,
) {
	return func(ctx context.Context) ([]internal.Service, error) {
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
}

func resolveService(conn grpc.ClientConnInterface) func(
	context.Context,
	internal.Service,
) (Service, error) {
	return func(ctx context.Context, d internal.Service) (Service, error) {
		stub := gr.NewServerReflectionClient(conn)
		c := grpcreflect.NewClient(ctx, stub)
		descriptor, err := c.ResolveService(d.Unwrap())
		if err != nil {
			switch {
			case isGrpcErr(err):
				st, _ := status.FromError(err)
				err = fmt.Errorf("resolve service %s: %w", d.Unwrap(), st.Err())
				return nil, err
			case isElementNotFoundErr(err):
				err := &internal.NotFound{Type: "service", Ident: d.Unwrap()}
				return nil, fmt.Errorf("resolve service: %w", err)
			case isProtocolError(err):
				err := fmt.Errorf("resolve service %s: %w", d.Unwrap(), err)
				return nil, err
			default:
				// Should never happen as all errors should be caught by one
				// of the above options
				return nil, errdefs.InternalWithMsg("resolve service: %v", err)
			}
		}
		if err != nil {
			return nil, err
		}
		return newService(descriptor)
	}
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
