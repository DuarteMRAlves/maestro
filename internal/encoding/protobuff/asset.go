package protobuff

import (
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/assert"
	"github.com/DuarteMRAlves/maestro/internal/asset"
)

func MarshalAsset(a *asset.Asset) *pb.Asset {
	return &pb.Asset{Id: MarshalID(a.Id), Name: a.Name, Image: a.Image}
}

func UnmarshalAsset(p *pb.Asset) (*asset.Asset, error) {
	if ok, err := assert.ArgNotNil(p, "p"); !ok {
		return nil, err
	}
	return &asset.Asset{
		Id:    UnmarshalId(p.Id),
		Name:  p.Name,
		Image: p.Image,
	}, nil
}
