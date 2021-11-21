package protobuff

import (
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/assert"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/stage"
)

func MarshalStage(s *stage.Stage) *pb.Stage {
	return &pb.Stage{
		Name:    s.Name,
		Asset:   s.Asset,
		Service: s.Service,
		Method:  s.Method,
	}
}

func UnmarshalStage(p *pb.Stage) (*stage.Stage, error) {
	if ok, err := assert.ArgNotNil(p, "p"); !ok {
		return nil, errdefs.InternalWithError(err)
	}
	return &stage.Stage{
		Name:    p.Name,
		Asset:   p.Asset,
		Service: p.Service,
		Method:  p.Method,
	}, nil
}
