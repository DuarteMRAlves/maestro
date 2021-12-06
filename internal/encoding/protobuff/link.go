package protobuff

import (
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/link"
	"github.com/DuarteMRAlves/maestro/internal/validate"
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
	if ok, err := validate.ArgNotNil(p, "p"); !ok {
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
