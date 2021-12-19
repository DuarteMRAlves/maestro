package mock

import (
	"context"
	"fmt"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

// LinkManagementServer is a mocking of a link management server to be used
// during tests. By default, methods are not implemented and raise an error.
//
// The server can be initialized with a specific function for each grpc method
// that will be called.
type LinkManagementServer struct {
	pb.UnimplementedLinkManagementServer
	CreateLinkFn func(ctx context.Context, config *pb.Link) (
		*emptypb.Empty,
		error,
	)
	GetLinkFn func(
		query *pb.Link,
		stream pb.LinkManagement_GetServer,
	) error
}

func (s *LinkManagementServer) Create(
	ctx context.Context,
	config *pb.Link,
) (*emptypb.Empty, error) {
	if s.CreateLinkFn != nil {
		return s.CreateLinkFn(ctx, config)
	}
	return &emptypb.Empty{}, fmt.Errorf(
		"method CreateLink not configured but called with config %v",
		config)
}

func (s *LinkManagementServer) Get(
	query *pb.Link,
	stream pb.LinkManagement_GetServer,
) error {
	if s.GetLinkFn != nil {
		return s.GetLinkFn(query, stream)
	}
	return fmt.Errorf(
		"method GetLink not configured but called with query %v",
		query)
}
