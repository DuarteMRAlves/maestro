package protobuff

import (
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/validate"
)

// MarshalAsset returns a new protobuf representations of the received asset.
func MarshalAsset(a *api.Asset) (*pb.Asset, error) {
	if ok, err := validate.ArgNotNil(a, "a"); !ok {
		return nil, err
	}
	return &pb.Asset{Name: string(a.Name), Image: a.Image}, nil
}

// UnmarshalAsset returns a new asset from its protobuf representation.
func UnmarshalAsset(p *pb.Asset) (*api.Asset, error) {
	if ok, err := validate.ArgNotNil(p, "p"); !ok {
		return nil, err
	}
	return &api.Asset{
		Name:  api.AssetName(p.Name),
		Image: p.Image,
	}, nil
}
