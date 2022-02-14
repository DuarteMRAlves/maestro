package client

import (
	"context"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/api"
)

func (c *client) StartExecution(
	ctx context.Context,
	req *api.StartExecutionRequest,
) error {
	pbReq := &pb.StartExecutionRequest{
		Orchestration: string(req.Orchestration),
	}

	stub := pb.NewExecutionManagementClient(c.conn)

	_, err := stub.Start(ctx, pbReq)

	return ErrorFromGrpcError(err)
}
