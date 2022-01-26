package client

import (
	"context"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"io"
	"time"
)

func (c *client) GetAsset(
	ctx context.Context,
	req *pb.GetAssetRequest,
) ([]*pb.Asset, error) {
	stub := pb.NewAssetManagementClient(c.conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	stream, err := stub.Get(ctx, req)
	if err != nil {
		return nil, ErrorFromGrpcError(err)
	}
	assets := make([]*pb.Asset, 0)
	for {
		a, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, ErrorFromGrpcError(err)
		}
		assets = append(assets, a)
	}
	return assets, nil
}

func (c *client) GetStage(
	ctx context.Context,
	req *pb.GetStageRequest,
) ([]*pb.Stage, error) {

	stub := pb.NewStageManagementClient(c.conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	stream, err := stub.Get(ctx, req)
	if err != nil {
		return nil, ErrorFromGrpcError(err)
	}
	stages := make([]*pb.Stage, 0)
	for {
		s, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, ErrorFromGrpcError(err)
		}
		stages = append(stages, s)
	}
	return stages, nil
}

func (c *client) GetLink(
	ctx context.Context,
	query *pb.Link,
) ([]*pb.Link, error) {

	stub := pb.NewLinkManagementClient(c.conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	stream, err := stub.Get(ctx, query)
	if err != nil {
		return nil, ErrorFromGrpcError(err)
	}
	links := make([]*pb.Link, 0)
	for {
		l, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, ErrorFromGrpcError(err)
		}
		links = append(links, l)
	}
	return links, nil
}

func (c *client) GetOrchestration(
	ctx context.Context,
	req *pb.GetOrchestrationRequest,
) ([]*pb.Orchestration, error) {

	stub := pb.NewOrchestrationManagementClient(c.conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	stream, err := stub.Get(ctx, req)
	if err != nil {
		return nil, ErrorFromGrpcError(err)
	}
	orchestrations := make([]*pb.Orchestration, 0)
	for {
		o, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, ErrorFromGrpcError(err)
		}
		orchestrations = append(orchestrations, o)
	}
	return orchestrations, nil
}
