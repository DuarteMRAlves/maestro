package protobuff

import (
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/asset"
	"github.com/DuarteMRAlves/maestro/internal/identifier"
)

func MarshalAsset(a *asset.Asset) *pb.Asset {
	return &pb.Asset{Id: MarshalID(a.Id), Name: a.Name, Image: a.Image}
}

func MarshalID(id identifier.Id) *pb.Id {
	return &pb.Id{Val: id.Val}
}
