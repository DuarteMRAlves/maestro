package protobuff

import (
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/asset"
	"github.com/DuarteMRAlves/maestro/internal/validate"
)

func MarshalAsset(a *asset.Asset) *pb.Asset {
	return &pb.Asset{Name: a.Name, Image: a.Image}
}

func UnmarshalAsset(p *pb.Asset) (*asset.Asset, error) {
	if ok, err := validate.ArgNotNil(p, "p"); !ok {
		return nil, err
	}
	return &asset.Asset{
		Name:  p.Name,
		Image: p.Image,
	}, nil
}
