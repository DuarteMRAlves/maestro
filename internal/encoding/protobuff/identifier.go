package protobuff

import (
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/identifier"
)

func MarshalID(id identifier.Id) *pb.Id {
	return &pb.Id{Val: id.Val}
}

func UnmarshalId(p *pb.Id) identifier.Id {
	if p == nil {
		return identifier.Empty()
	}
	return identifier.Id{Val: p.Val}
}
