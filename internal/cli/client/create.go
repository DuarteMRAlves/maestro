package client

import (
	"context"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/cli/resources"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
)

// CreateResource creates a resource of any kind.
// It does no checking of any pre-condition for the resource creation.
func CreateResource(
	ctx context.Context,
	resource *resources.Resource,
	addr string,
) error {
	switch {
	case resources.IsAssetKind(resource):
		a := &resources.AssetResource{}
		if err := resources.MarshalResource(a, resource); err != nil {
			return err
		}
		if err := CreateAsset(ctx, a, addr); err != nil {
			return err
		}
		return nil
	case resources.IsStageKind(resource):
		s := &resources.StageResource{}
		if err := resources.MarshalResource(s, resource); err != nil {
			return err
		}
		if err := CreateStage(ctx, s, addr); err != nil {
			return err
		}
		return nil
	case resources.IsLinkKind(resource):
		l := &resources.LinkResource{}
		if err := resources.MarshalResource(l, resource); err != nil {
			return err
		}
		if err := CreateLink(ctx, l, addr); err != nil {
			return err
		}
		return nil
	}
	return errdefs.InvalidArgumentWithMsg("unknown kind %v", resource.Kind)
}

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

	return ErrorFromGrpcError(err)
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

	return ErrorFromGrpcError(err)
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

	return ErrorFromGrpcError(err)
}
