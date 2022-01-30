package pb

import (
	"context"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/api"
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
	pbReq *pb.CreateAssetRequest,
) (*emptypb.Empty, error) {

	var (
		req api.CreateAssetRequest
	)
	var err error
	var grpcErr error = nil

	UnmarshalCreateAssetRequest(&req, pbReq)
	err = s.api.CreateAsset(&req)
	if err != nil {
		grpcErr = GrpcErrorFromError(err)
	}
	return &emptypb.Empty{}, grpcErr
}

func (s *assetManagementServer) Get(
	pbQuery *pb.GetAssetRequest,
	stream pb.AssetManagement_GetServer,
) error {

	var (
		query api.GetAssetRequest
		err   error
	)

	UnmarshalGetAssetRequest(&query, pbQuery)

	assets, err := s.api.GetAsset(&query)
	if err != nil {
		return GrpcErrorFromError(err)
	}
	for _, a := range assets {
		pbAsset, err := MarshalAsset(a)
		if err != nil {
			return err
		}
		stream.Send(pbAsset)
	}
	return nil
}
