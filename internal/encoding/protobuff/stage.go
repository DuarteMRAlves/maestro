package protobuff

import (
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/stage"
	"github.com/DuarteMRAlves/maestro/internal/validate"
)

// MarshalStage creates a protobuf message representing the Stage from the Stage
// structure.
func MarshalStage(s *stage.Stage) *pb.Stage {
	return &pb.Stage{
		Name:    s.Name,
		Asset:   s.Asset,
		Service: s.Service,
		Method:  s.Method,
		Address: s.Address,
	}
}

// UnmarshalStage creates a Stage struct from a protobuf message representing
// the stage.
func UnmarshalStage(p *pb.Stage) (*stage.Stage, error) {
	if ok, err := validate.ArgNotNil(p, "p"); !ok {
		return nil, errdefs.InvalidArgumentWithError(err)
	}
	return &stage.Stage{
		Name:    p.Name,
		Asset:   p.Asset,
		Service: p.Service,
		Method:  p.Method,
		Address: p.Address,
	}, nil
}
