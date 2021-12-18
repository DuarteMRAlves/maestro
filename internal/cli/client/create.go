package client

import (
	"context"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/cli/resources"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
)

// CreateResource creates a resource of any kind.
// It does no checking of any pre-condition for the resource creation.
func (c *client) CreateResource(
	ctx context.Context,
	resource *resources.Resource,
) error {
	switch {
	case resources.IsAssetKind(resource):
		a, ok := resource.Spec.(*resources.AssetSpec)
		if !ok {
			return errdefs.InternalWithMsg("asset spec cast failed: %v", a)
		}
		if err := c.CreateAsset(ctx, a); err != nil {
			return err
		}
		return nil
	case resources.IsStageKind(resource):
		s, ok := resource.Spec.(*resources.StageSpec)
		if !ok {
			return errdefs.InternalWithMsg("stage spec cast failed: %v", s)
		}
		if err := c.CreateStage(ctx, s); err != nil {
			return err
		}
		return nil
	case resources.IsLinkKind(resource):
		l, ok := resource.Spec.(*resources.LinkSpec)
		if !ok {
			return errdefs.InternalWithMsg("link spec cast failed: %v", l)
		}
		if err := c.CreateLink(ctx, l); err != nil {
			return err
		}
		return nil
	}
	return errdefs.InvalidArgumentWithMsg("unknown kind %v", resource.Kind)
}

func (c *client) CreateAsset(
	ctx context.Context,
	asset *resources.AssetSpec,
) error {
	a := &pb.Asset{
		Name:  asset.Name,
		Image: asset.Image,
	}

	stub := pb.NewAssetManagementClient(c.conn)

	_, err := stub.Create(ctx, a)

	return ErrorFromGrpcError(err)
}

func (c *client) CreateStage(
	ctx context.Context,
	stage *resources.StageSpec,
) error {
	s := &pb.Stage{
		Name:    stage.Name,
		Asset:   stage.Asset,
		Service: stage.Service,
		Method:  stage.Method,
		Address: stage.Address,
	}

	stub := pb.NewStageManagementClient(c.conn)

	_, err := stub.Create(ctx, s)

	return ErrorFromGrpcError(err)
}

func (c *client) CreateLink(
	ctx context.Context,
	link *resources.LinkSpec,
) error {
	l := &pb.Link{
		Name:        link.Name,
		SourceStage: link.SourceStage,
		SourceField: link.SourceField,
		TargetStage: link.TargetStage,
		TargetField: link.TargetField,
	}

	stub := pb.NewLinkManagementClient(c.conn)

	_, err := stub.Create(ctx, l)

	return ErrorFromGrpcError(err)
}

func (c *client) CreateOrchestration(
	ctx context.Context,
	orchestration *resources.OrchestrationSpec,
) error {
	o := &pb.Orchestration{
		Name:  orchestration.Name,
		Links: orchestration.Links,
	}

	stub := pb.NewOrchestrationManagementClient(c.conn)

	_, err := stub.Create(ctx, o)

	return ErrorFromGrpcError(err)
}
