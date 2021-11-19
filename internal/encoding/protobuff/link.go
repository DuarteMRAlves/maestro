package protobuff

import (
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/assert"
	"github.com/DuarteMRAlves/maestro/internal/link"
)

func MarshalLink(l *link.Link) *pb.Link {
	return &pb.Link{
		Name:        l.Name,
		SourceStage: l.SourceStage,
		SourceField: l.SourceField,
		TargetStage: l.TargetStage,
		TargetField: l.TargetField,
	}
}

func UnmarshalLink(p *pb.Link) (*link.Link, error) {
	if ok, err := assert.ArgNotNil(p, "p"); !ok {
		return nil, err
	}
	return &link.Link{
		Name:        p.Name,
		SourceStage: p.SourceStage,
		SourceField: p.SourceField,
		TargetStage: p.TargetStage,
		TargetField: p.TargetField,
	}, nil
}
