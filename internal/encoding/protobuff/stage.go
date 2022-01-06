package protobuff

import (
	"github.com/DuarteMRAlves/maestro/api/pb"
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/validate"
)

// MarshalStage creates a protobuf message representing the Stage from the Stage
// structure.
func MarshalStage(s *apitypes.Stage) (*pb.Stage, error) {
	if ok, err := validate.ArgNotNil(s, "s"); !ok {
		return nil, err
	}
	pbStage := &pb.Stage{
		Name:    s.Name,
		Asset:   string(s.Asset),
		Service: s.Service,
		Rpc:     s.Rpc,
		Address: s.Address,
		Host:    s.Host,
		Port:    s.Port,
	}
	return pbStage, nil
}

// UnmarshalStage creates a Stage struct from a protobuf message representing
// the stage.
func UnmarshalStage(p *pb.Stage) (*apitypes.Stage, error) {
	if ok, err := validate.ArgNotNil(p, "p"); !ok {
		return nil, errdefs.InvalidArgumentWithError(err)
	}
	return &apitypes.Stage{
		Name:    p.Name,
		Asset:   apitypes.AssetName(p.Asset),
		Service: p.Service,
		Rpc:     p.Rpc,
		Address: p.Address,
		Host:    p.Host,
		Port:    p.Port,
	}, nil
}
