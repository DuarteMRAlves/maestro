package protobuff

import (
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/assert"
	"github.com/DuarteMRAlves/maestro/internal/blueprint"
)

func MarshalStage(s *blueprint.Stage) *pb.Stage {
	return &pb.Stage{
		Id:      MarshalID(s.Id),
		Name:    s.Name,
		AssetId: MarshalID(s.AssetId),
		Service: s.Service,
		Method:  s.Method,
	}
}

func UnmarshalStage(p *pb.Stage) (*blueprint.Stage, error) {
	if ok, err := assert.ArgNotNil(p, "p"); !ok {
		return nil, err
	}
	return &blueprint.Stage{
		Id:      UnmarshalId(p.Id),
		Name:    p.Name,
		AssetId: UnmarshalId(p.AssetId),
		Service: p.Service,
		Method:  p.Method,
	}, nil
}
