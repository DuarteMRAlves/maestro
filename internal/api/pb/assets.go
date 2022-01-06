package pb

import (
	"context"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/api"
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"github.com/DuarteMRAlves/maestro/internal/encoding/protobuff"
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

	var a *apitypes.Asset
	var err error
	var grpcErr error = nil

	if a, err = protobuff.UnmarshalAsset(pbAsset); err != nil {
		return &emptypb.Empty{}, err
	}
	err = s.api.CreateAsset(a)
	if err != nil {
		grpcErr = GrpcErrorFromError(err)
	}
	return &emptypb.Empty{}, grpcErr
}

func (s *assetManagementServer) Get(
	pbQuery *pb.Asset,
	stream pb.AssetManagement_GetServer,
) error {

	var query *apitypes.Asset
	var err error

	if query, err = protobuff.UnmarshalAsset(pbQuery); err != nil {
		return err
	}

	assets := s.api.GetAsset(query)
	for _, a := range assets {
		pbAsset, err := protobuff.MarshalAsset(a)
		if err != nil {
			return err
		}
		stream.Send(pbAsset)
	}
	return nil
}
