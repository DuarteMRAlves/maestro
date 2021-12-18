package client

import (
	"context"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/cli/resources"
	"google.golang.org/grpc"
)

// Client offers an interface to connect to the maestro server
type Client interface {
	CreateResource(ctx context.Context, r *resources.Resource) error
	CreateAsset(ctx context.Context, a *resources.AssetSpec) error
	CreateStage(ctx context.Context, s *resources.StageSpec) error
	CreateLink(ctx context.Context, l *resources.LinkSpec) error
	CreateOrchestration(
		ctx context.Context,
		o *resources.OrchestrationSpec,
	) error

	GetAsset(ctx context.Context, query *pb.Asset) ([]*pb.Asset, error)
	GetStage(ctx context.Context, query *pb.Stage) ([]*pb.Stage, error)
	GetLink(ctx context.Context, query *pb.Link) ([]*pb.Link, error)
}

type client struct {
	conn grpc.ClientConnInterface
}

func New(conn grpc.ClientConnInterface) Client {
	return &client{conn: conn}
}
