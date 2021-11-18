package protobuff

import (
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/assert"
	"github.com/DuarteMRAlves/maestro/internal/blueprint"
)

func MarshalLink(l *blueprint.Link) *pb.Link {
	return &pb.Link{
		SourceStage: l.SourceStage,
		SourceField: l.SourceField,
		TargetStage: l.TargetStage,
		TargetField: l.TargetField,
	}
}

func UnmarshalLink(p *pb.Link) (*blueprint.Link, error) {
	if ok, err := assert.ArgNotNil(p, "p"); !ok {
		return nil, err
	}
	return &blueprint.Link{
		SourceStage: p.SourceStage,
		SourceField: p.SourceField,
		TargetStage: p.TargetStage,
		TargetField: p.TargetField,
	}, nil
}
