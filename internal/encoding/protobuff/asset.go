package protobuff

import (
	"github.com/DuarteMRAlves/maestro/api/pb"
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"github.com/DuarteMRAlves/maestro/internal/validate"
)

// MarshalAsset returns a new protobuf representations of the received asset.
func MarshalAsset(a *apitypes.Asset) (*pb.Asset, error) {
	if ok, err := validate.ArgNotNil(a, "a"); !ok {
		return nil, err
	}
	return &pb.Asset{Name: a.Name, Image: a.Image}, nil
}

// UnmarshalAsset returns a new asset from its protobuf representation.
func UnmarshalAsset(p *pb.Asset) (*apitypes.Asset, error) {
	if ok, err := validate.ArgNotNil(p, "p"); !ok {
		return nil, err
	}
	return &apitypes.Asset{
		Name:  p.Name,
		Image: p.Image,
	}, nil
}
