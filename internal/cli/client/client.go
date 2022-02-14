package client

import (
	"context"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/api"
	"google.golang.org/grpc"
)

// Client offers an interface to connect to the maestro server
type Client interface {
	CreateAsset(context.Context, *api.CreateAssetRequest) error
	CreateStage(context.Context, *api.CreateStageRequest) error
	CreateLink(context.Context, *api.CreateLinkRequest) error
	CreateOrchestration(context.Context, *api.CreateOrchestrationRequest) error

	GetAsset(context.Context, *pb.GetAssetRequest) ([]*pb.Asset, error)
	GetStage(context.Context, *pb.GetStageRequest) ([]*pb.Stage, error)
	GetLink(context.Context, *pb.GetLinkRequest) ([]*pb.Link, error)
	GetOrchestration(context.Context, *pb.GetOrchestrationRequest) (
		[]*pb.Orchestration,
		error,
	)

	StartExecution(context.Context, *api.StartExecutionRequest) error
}

type client struct {
	conn grpc.ClientConnInterface
}

func New(conn grpc.ClientConnInterface) Client {
	return &client{conn: conn}
}
