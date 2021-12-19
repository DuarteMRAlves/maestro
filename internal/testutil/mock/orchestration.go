package mock

import (
	"context"
	"fmt"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

// OrchestrationManagementServer is a mocking of an orchestration management
// server to be used during tests. By default, methods are not implemented and
// raise an error.
//
// The server can be initialized with a specific function for each grpc method
// that will be called.
type OrchestrationManagementServer struct {
	pb.UnimplementedOrchestrationManagementServer
	CreateOrchestrationFn func(ctx context.Context, config *pb.Orchestration) (
		*emptypb.Empty,
		error,
	)
	GetOrchestrationFn func(
		query *pb.Orchestration,
		stream pb.OrchestrationManagement_GetServer,
	) error
}

func (s *OrchestrationManagementServer) Create(
	ctx context.Context,
	config *pb.Orchestration,
) (*emptypb.Empty, error) {
	if s.CreateOrchestrationFn != nil {
		return s.CreateOrchestrationFn(ctx, config)
	}
	return &emptypb.Empty{}, fmt.Errorf(
		"method CreateOrchestration not configured but called with config %v",
		config)
}

func (s *OrchestrationManagementServer) Get(
	query *pb.Orchestration,
	stream pb.OrchestrationManagement_GetServer,
) error {
	if s.GetOrchestrationFn != nil {
		return s.GetOrchestrationFn(query, stream)
	}
	return fmt.Errorf(
		"method GetOrchestration not configured but called with query %v",
		query)
}
