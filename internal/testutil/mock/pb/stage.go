package pb

import (
	"context"
	"fmt"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

// StageManagementServer is a mocking of a stage management server to be used
// during tests. By default, methods are not implemented and raise an error.
//
// The server can be initialized with a specific function for each grpc method
// that will be called.
type StageManagementServer struct {
	pb.UnimplementedStageManagementServer
	CreateStageFn func(ctx context.Context, config *pb.Stage) (
		*emptypb.Empty,
		error,
	)
	GetStageFn func(
		query *pb.Stage,
		stream pb.StageManagement_GetServer,
	) error
}

func (s *StageManagementServer) Create(
	ctx context.Context,
	config *pb.Stage,
) (*emptypb.Empty, error) {
	if s.CreateStageFn != nil {
		return s.CreateStageFn(ctx, config)
	}
	return &emptypb.Empty{}, fmt.Errorf(
		"method CreateStage not configured but called with config %v",
		config)
}

func (s *StageManagementServer) Get(
	query *pb.Stage,
	stream pb.StageManagement_GetServer,
) error {
	if s.GetStageFn != nil {
		return s.GetStageFn(query, stream)
	}
	return fmt.Errorf(
		"method GetStage not configured but called with query %v",
		query)
}
