package protobuff

import (
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/assert"
	"github.com/DuarteMRAlves/maestro/internal/blueprint"
)

func MarshalLink(l *blueprint.Link) *pb.Link {
	return &pb.Link{
		SourceId:    MarshalID(l.SourceId),
		SourceField: l.SourceField,
		TargetId:    MarshalID(l.TargetId),
		TargetField: l.TargetField,
	}
}

func UnmarshalLink(p *pb.Link) (*blueprint.Link, error) {
	if ok, err := assert.ArgNotNil(p, "p"); !ok {
		return nil, err
	}
	return &blueprint.Link{
		SourceId:    UnmarshalId(p.SourceId),
		SourceField: p.SourceField,
		TargetId:    UnmarshalId(p.TargetId),
		TargetField: p.TargetField,
	}, nil
}
