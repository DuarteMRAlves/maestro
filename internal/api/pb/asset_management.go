package pb

import (
	"context"
	"errors"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/asset"
	"github.com/DuarteMRAlves/maestro/internal/encoding/protobuff"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type assetManagementServer struct {
	pb.UnimplementedAssetManagementServer
	api api.InternalAPI
}

func NewAssetManagementServer(api api.InternalAPI) pb.AssetManagementServer {
	return &assetManagementServer{api: api}
}

func (s *assetManagementServer) Create(
	ctx context.Context,
	pbAsset *pb.Asset,
) (*emptypb.Empty, error) {

	var a *asset.Asset
	var err error
	var grpcErr error = nil

	if a, err = protobuff.UnmarshalAsset(pbAsset); err != nil {
		return &emptypb.Empty{}, err
	}
	err = s.api.CreateAsset(a)
	if err != nil {
		var alreadyExists asset.AlreadyExists
		if errors.As(err, &alreadyExists) {
			grpcErr = status.Error(codes.AlreadyExists, alreadyExists.Error())
		} else {
			grpcErr = status.Error(codes.Unknown, err.Error())
		}
	}
	return &emptypb.Empty{}, grpcErr
}

func (s *assetManagementServer) Get(
	pbQuery *pb.Asset,
	stream pb.AssetManagement_GetServer,
) error {

	var query *asset.Asset
	var err error

	if query, err = protobuff.UnmarshalAsset(pbQuery); err != nil {
		return err
	}

	assets := s.api.GetAsset(query)
	for _, a := range assets {
		stream.Send(protobuff.MarshalAsset(a))
	}
	return nil
}
