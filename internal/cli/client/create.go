package client

import (
	"context"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/api"
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
	case resource.IsAssetKind():
		req, ok := resource.Spec.(*api.CreateAssetRequest)
		if !ok {
			return errdefs.InternalWithMsg(
				"create asset request spec cast failed: %v",
				req,
			)
		}
		if err := c.CreateAsset(ctx, req); err != nil {
			return err
		}
		return nil
	case resource.IsStageKind():
		s, ok := resource.Spec.(*api.CreateStageRequest)
		if !ok {
			return errdefs.InternalWithMsg("stage spec cast failed: %v", s)
		}
		if err := c.CreateStage(ctx, s); err != nil {
			return err
		}
		return nil
	case resource.IsLinkKind():
		l, ok := resource.Spec.(*api.CreateLinkRequest)
		if !ok {
			return errdefs.InternalWithMsg("link spec cast failed: %v", l)
		}
		if err := c.CreateLink(ctx, l); err != nil {
			return err
		}
		return nil
	case resource.IsOrchestrationKind():
		o, ok := resource.Spec.(*api.CreateOrchestrationRequest)
		if !ok {
			return errdefs.InternalWithMsg(
				"orchestration spec cast failed> %v",
				o,
			)
		}
		if err := c.CreateOrchestration(ctx, o); err != nil {
			return err
		}
		return nil
	default:
		return errdefs.InvalidArgumentWithMsg("unknown kind %v", resource.Kind)
	}
}

func (c *client) CreateAsset(
	ctx context.Context,
	req *api.CreateAssetRequest,
) error {
	a := &pb.CreateAssetRequest{
		Name:  string(req.Name),
		Image: req.Image,
	}

	stub := pb.NewAssetManagementClient(c.conn)

	_, err := stub.Create(ctx, a)

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

	stub := pb.NewStageManagementClient(c.conn)

	_, err := stub.Create(ctx, pbReq)

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

	stub := pb.NewLinkManagementClient(c.conn)

	_, err := stub.Create(ctx, l)

	return ErrorFromGrpcError(err)
}

func (c *client) CreateOrchestration(
	ctx context.Context,
	req *api.CreateOrchestrationRequest,
) error {
	links := make([]string, 0, len(req.Links))
	for _, l := range req.Links {
		links = append(links, string(l))
	}
	pbReq := &pb.CreateOrchestrationRequest{
		Name:  string(req.Name),
		Links: links,
	}

	stub := pb.NewOrchestrationManagementClient(c.conn)

	_, err := stub.Create(ctx, pbReq)

	return ErrorFromGrpcError(err)
}
