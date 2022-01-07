package protobuff

import (
	"github.com/DuarteMRAlves/maestro/api/pb"
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"github.com/DuarteMRAlves/maestro/internal/validate"
)

// MarshalLink returns a protobuf message for the given link.
func MarshalLink(l *apitypes.Link) (*pb.Link, error) {
	if ok, err := validate.ArgNotNil(l, "l"); !ok {
		return nil, err
	}
	pbLink := &pb.Link{
		Name:        string(l.Name),
		SourceStage: string(l.SourceStage),
		SourceField: l.SourceField,
		TargetStage: string(l.TargetStage),
		TargetField: l.TargetField,
	}
	return pbLink, nil
}

// UnmarshalLink returns the link represented by the given protobuf message.
func UnmarshalLink(p *pb.Link) (*apitypes.Link, error) {
	if ok, err := validate.ArgNotNil(p, "p"); !ok {
		return nil, err
	}
	return &apitypes.Link{
		Name:        apitypes.LinkName(p.Name),
		SourceStage: apitypes.StageName(p.SourceStage),
		SourceField: p.SourceField,
		TargetStage: apitypes.StageName(p.TargetStage),
		TargetField: p.TargetField,
	}, nil
}
