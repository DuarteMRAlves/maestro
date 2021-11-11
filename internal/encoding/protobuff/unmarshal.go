package protobuff

import (
	pb "github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/assert"
	"github.com/DuarteMRAlves/maestro/internal/asset"
	"github.com/DuarteMRAlves/maestro/internal/identifier"
)

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

func UnmarshalId(p *pb.Id) identifier.Id {
	if p == nil {
		return identifier.Empty()
	}
	return identifier.Id{Val: p.Val}
}
