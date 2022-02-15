package client

import (
	"context"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/api"
)

func (c *client) CreateAsset(
	ctx context.Context,
	req *api.CreateAssetRequest,
) error {
	a := &pb.CreateAssetRequest{
		Name:  string(req.Name),
		Image: req.Image,
	}

	stub := pb.NewArchitectureManagementClient(c.conn)

	_, err := stub.CreateAsset(ctx, a)

	return ErrorFromGrpcError(err)
}

func (c *client) CreateStage(
	ctx context.Context,
	req *api.CreateStageRequest,
) error {
	pbReq := &pb.CreateStageRequest{
		Name:    string(req.Name),
		Asset:   string(req.Asset),
		Service: req.Service,
		Rpc:     req.Rpc,
		Address: req.Address,
		Host:    req.Host,
		Port:    req.Port,
	}

	stub := pb.NewArchitectureManagementClient(c.conn)

	_, err := stub.CreateStage(ctx, pbReq)

	return ErrorFromGrpcError(err)
}

func (c *client) CreateLink(
	ctx context.Context,
	link *api.CreateLinkRequest,
) error {
	l := &pb.CreateLinkRequest{
		Name:        string(link.Name),
		SourceStage: string(link.SourceStage),
		SourceField: link.SourceField,
		TargetStage: string(link.TargetStage),
		TargetField: link.TargetField,
	}

	stub := pb.NewArchitectureManagementClient(c.conn)

	_, err := stub.CreateLink(ctx, l)

	return ErrorFromGrpcError(err)
}

func (c *client) CreateOrchestration(
	ctx context.Context,
	req *api.CreateOrchestrationRequest,
) error {
	pbReq := &pb.CreateOrchestrationRequest{
		Name: string(req.Name),
	}

	stub := pb.NewArchitectureManagementClient(c.conn)

	_, err := stub.CreateOrchestration(ctx, pbReq)

	return ErrorFromGrpcError(err)
}
