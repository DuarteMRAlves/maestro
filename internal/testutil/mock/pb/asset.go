package pb

import (
	"context"
	"fmt"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

// AssetManagementServer is a mocking of an asset management server to be used
// during tests. By default, methods are not implemented and raise an error.
//
// The server can be initialized with a specific function for each grpc method
// that will be called.
type AssetManagementServer struct {
	pb.UnimplementedAssetManagementServer
	CreateAssetFn func(ctx context.Context, config *pb.Asset) (
		*emptypb.Empty,
		error,
	)
	GetAssetFn func(
		query *pb.Asset,
		stream pb.AssetManagement_GetServer,
	) error
}

func (s *AssetManagementServer) Create(
	ctx context.Context,
	config *pb.Asset,
) (*emptypb.Empty, error) {
	if s.CreateAssetFn != nil {
		return s.CreateAssetFn(ctx, config)
	}
	return &emptypb.Empty{}, fmt.Errorf(
		"method CreateAsset not configured but called with config %v",
		config)
}

func (s *AssetManagementServer) Get(
	query *pb.Asset,
	stream pb.AssetManagement_GetServer,
) error {
	if s.GetAssetFn != nil {
		return s.GetAssetFn(query, stream)
	}
	return fmt.Errorf(
		"method GetAsset not configured but called with query %v",
		query)
}
