package invoke

import (
	"context"
	"github.com/DuarteMRAlves/maestro/internal/domain"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/jhump/protoreflect/grpcreflect"
	"google.golang.org/grpc"
	gr "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

const reflectionServiceName = "grpc.reflection.v1alpha.ServerReflection"

func listServices(conn grpc.ClientConnInterface) func(context.Context) (
	[]domain.Service,
	error,
) {
	return func(ctx context.Context) ([]domain.Service, error) {
		stub := gr.NewServerReflectionClient(conn)
		c := grpcreflect.NewClient(ctx, stub)
		all, err := c.ListServices()
		if err != nil {
			return nil, handleGrpcError(err, "list services: ")
		}
		// Filter the reflection service
		services := make([]domain.Service, 0, len(all)-1)
		for _, s := range all {
			if s != reflectionServiceName {
				name, err := domain.NewService(s)
				if err != nil {
					return nil, err
				}
				services = append(services, name)
			}
		}
		return services, nil
	}
}

func resolveService(conn grpc.ClientConnInterface) func(
	context.Context,
	domain.Service,
) (Service, error) {
	return func(ctx context.Context, d domain.Service) (Service, error) {
		stub := gr.NewServerReflectionClient(conn)
		c := grpcreflect.NewClient(ctx, stub)
		descriptor, err := c.ResolveService(d.Unwrap())
		if err != nil {
			switch {
			case isGrpcErr(err):
				return nil, handleGrpcError(err, "resolve service: ")
			case isElementNotFoundErr(err):
				return nil, errdefs.NotFoundWithMsg(
					"resolve service: %v",
					err.Error(),
				)
			case isProtocolError(err):
				return nil, errdefs.UnknownWithError(err)
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