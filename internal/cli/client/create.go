package client

import (
	"context"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/cli/resources"
)

func CreateAsset(
	ctx context.Context,
	asset *resources.AssetResource,
	addr string,
) error {
	a := &pb.Asset{
		Name:  asset.Name,
		Image: asset.Image,
	}
	conn := NewConnection(addr)
	defer conn.Close()

	c := pb.NewAssetManagementClient(conn)

	_, err := c.Create(ctx, a)

	return err
}

func CreateStage(
	ctx context.Context,
	stage *resources.StageResource,
	addr string,
) error {
	s := &pb.Stage{
		Name:    stage.Name,
		Asset:   stage.Asset,
		Service: stage.Service,
		Method:  stage.Method,
	}
	conn := NewConnection(addr)
	defer conn.Close()

	c := pb.NewStageManagementClient(conn)

	_, err := c.Create(ctx, s)

	return err
}

func CreateLink(
	ctx context.Context,
	link *resources.LinkResource,
	addr string,
) error {
	l := &pb.Link{
		Name:        link.Name,
		SourceStage: link.SourceStage,
		SourceField: link.SourceField,
		TargetStage: link.TargetStage,
		TargetField: link.TargetField,
	}
	conn := NewConnection(addr)
	defer conn.Close()

	c := pb.NewLinkManagementClient(conn)
	
	_, err := c.Create(ctx, l)

	return err
}
