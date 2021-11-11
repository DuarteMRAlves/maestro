package pb

import (
	"context"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/asset"
	"github.com/DuarteMRAlves/maestro/internal/encoding/protobuff"
	"github.com/DuarteMRAlves/maestro/internal/identifier"
)

var emptyIdPb = protobuff.MarshalID(identifier.Empty())

type assetManagementServer struct {
	pb.UnimplementedAssetManagementServer
	api api.InternalAPI
}

func NewAssetManagementServer(api api.InternalAPI) pb.AssetManagementServer {
	return &assetManagementServer{api: api}
}

func (s *assetManagementServer) Create(
	ctx context.Context,
	pbAsset *pb.Asset) (*pb.Id, error) {

	var a *asset.Asset
	var id identifier.Id
	var err error

	if a, err = protobuff.UnmarshalAsset(pbAsset); err != nil {
		return emptyIdPb, err
	}
	if id, err = s.api.Create(a); err != nil {
		return emptyIdPb, err
	}
	return protobuff.MarshalID(id), nil
}

func (s *assetManagementServer) List(
	_ *pb.SearchQuery,
	stream pb.AssetManagement_ListServer) error {

	assets, err := s.api.List()
	if err != nil {
		return err
	}
	for _, a := range assets {
		stream.Send(protobuff.MarshalAsset(a))
	}
	return nil
}
