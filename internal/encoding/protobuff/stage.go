package protobuff

import (
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/util"
)

// MarshalStage creates a protobuf message representing the Stage from the Stage
// structure.
func MarshalStage(s *api.Stage) (*pb.Stage, error) {
	if ok, err := util.ArgNotNil(s, "s"); !ok {
		return nil, err
	}
	pbStage := &pb.Stage{
		Name:    string(s.Name),
		Phase:   string(s.Phase),
		Asset:   string(s.Asset),
		Service: s.Service,
		Rpc:     s.Rpc,
		Address: s.Address,
	}
	return pbStage, nil
}

// UnmarshalStage creates a Stage struct from a protobuf message representing
// the stage.
func UnmarshalStage(p *pb.Stage) (*api.Stage, error) {
	if ok, err := util.ArgNotNil(p, "p"); !ok {
		return nil, errdefs.InvalidArgumentWithError(err)
	}
	return &api.Stage{
		Name:    api.StageName(p.Name),
		Phase:   api.StagePhase(p.Phase),
		Asset:   api.AssetName(p.Asset),
		Service: p.Service,
		Rpc:     p.Rpc,
		Address: p.Address,
	}, nil
}
